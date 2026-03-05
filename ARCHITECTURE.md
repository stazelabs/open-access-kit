# OAK (Open Access Kit) — Design Document

## Context

OAK is a portable, offline-first collection of privacy tools, secure communication apps, and curated knowledge designed to fit on a single USB thumbdrive. The project needs tooling to curate, download, verify, stage, and package quarterly releases (e.g., Q126 = Q1 2026) across media size tiers. A single command — `oak build` — drives the entire quarterly update cycle: detect latest upstream versions, download everything, verify signatures, assemble the image, and package it.

### Problem

Privacy tools like Tor Browser and Tails OS are critical for people facing censorship or surveillance, but downloading them requires internet access — exactly what may be restricted. OAK solves this by pre-packaging these tools onto removable media that can be physically shared and replicated.

### Goals

- Self-contained Go CLI (`oak`) that anyone can run to build an OAK image
- Quarterly release cadence with automated version detection
- Tiered images for different media sizes
- Offline-usable companion documentation and website
- Cryptographic verification of all included software
- GPG signing of the OAK image itself

---

## 1. Repo Layout

```
open-access-kit/
├── README.md                     # GitHub landing page
├── LICENSE
├── ARCHITECTURE.md               # This document
├── go.mod / go.sum
├── Makefile                      # build, test, lint, release targets
├── oak.yaml                      # Tier + source configuration
│
├── cmd/oak/
│   └── main.go                   # CLI entrypoint (cobra)
│
├── internal/
│   ├── config/                   # YAML config loading, structs
│   ├── source/                   # Source interface + implementations
│   │   ├── source.go             # Interface definition
│   │   ├── rsync.go              # Tor Browser, Tails
│   │   ├── git.go                # Onion sites directory
│   │   ├── http.go               # HTTP fetcher (version detection)
│   │   └── local.go              # Bundled local content
│   ├── version/                  # Version detection (scrape, RSS)
│   ├── verify/                   # GPG + checksum verification
│   ├── tier/                     # Size budgets, content selection
│   ├── stage/                    # Mirror -> image layout
│   ├── annotate/                 # README, MANIFEST, VERSION generation
│   ├── packaging/                # ZIP creation + GPG signing
│   ├── site/                     # Go-native Markdown->HTML renderer
│   └── pipeline/                 # Orchestrates full build
│
├── cli/                          # Cobra command definitions
│   ├── build.go                  # `oak build` (full pipeline)
│   ├── download.go               # `oak download`
│   ├── verify.go                 # `oak verify`
│   ├── stage.go                  # `oak stage`
│   ├── status.go                 # `oak status`
│   └── version.go                # `oak version`
│
├── content/                      # Educational content (Markdown source)
│   ├── guides/
│   │   ├── what-is-tor.md
│   │   ├── using-tor-browser.md
│   │   ├── what-is-tails.md
│   │   ├── privacy-basics.md
│   │   └── censorship-circumvention.md
│   └── templates/
│       ├── README.md.tmpl        # Root README template for images
│       └── site/                 # HTML templates for companion website
│
├── keys/                         # Upstream GPG public keys (checked in)
│   ├── torproject-signing.gpg
│   └── tails-signing.gpg
│
├── scripts/                      # Legacy bash scripts (reference only)
│
└── .github/workflows/ci.yml      # Go build + test + lint
```

**Runtime artifacts** (gitignored):
- `mirror/` — raw downloaded content
- `image/` — staged image layout
- `dist/` — packaged ZIP archives

---

## 2. Go CLI Architecture

### Framework

[spf13/cobra](https://github.com/spf13/cobra) — the standard for Go CLIs.

### Commands

| Command | Description |
|---------|-------------|
| `oak build --tier 64` | Full pipeline: download -> verify -> stage -> annotate -> package -> sign |
| `oak download [source]` | Fetch/mirror all or a specific source |
| `oak verify [source]` | Check GPG sigs + checksums of mirrored content |
| `oak stage --tier 64` | Build image directory from mirror |
| `oak status` | Show mirror state, sizes, detected versions |
| `oak version` | Print CLI version |

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config, -c` | `./oak.yaml` | Path to configuration file |
| `--tier, -t` | `64` | Target tier: `16`, `32`, `64`, `max` |
| `--mirror-dir` | `./mirror` | Path to mirror directory |
| `--image-dir` | `./image` | Path to image output directory |
| `--release` | auto (e.g., `Q126`) | Release name override |
| `--dry-run` | `false` | Show what would happen without doing it |
| `--verbose, -v` | `false` | Verbose output |

### Source Interface

```go
type Source interface {
    Name() string
    DetectVersion(ctx context.Context) (string, error)
    Download(ctx context.Context, mirrorDir string) error
    Verify(ctx context.Context, mirrorDir string) error
    Size(mirrorDir string) (int64, error)
    Stage(ctx context.Context, mirrorDir, imageDir string, tier TierConfig) error
}
```

Implementations:
- **RsyncSource** — Tor Browser, Tails (rsync protocol)
- **GitSource** — onion sites directory (git clone/pull)
- **HTTPSource** — generic HTTP downloads
- **LocalSource** — bundled educational content

Each is config-driven. Adding a new source means adding a YAML block, not new Go code (if it fits an existing type).

---

## 3. Configuration (`oak.yaml`)

```yaml
release: Q126

paths:
  mirror: ./mirror
  image: ./image
  output: ./dist

tiers:
  16:
    label: "16GB"
    budget_gb: 14          # Usable space after filesystem overhead
    sources:
      - tor-browser
      - onion-sites
      - educational-content
  32:
    label: "32GB"
    budget_gb: 29
    sources:
      - tor-browser
      - onion-sites
      - educational-content
      - tails
  64:
    label: "64GB"
    budget_gb: 58
    sources:
      - tor-browser
      - onion-sites
      - educational-content
      - tails
  max:
    label: "Max"
    budget_gb: 0           # No limit
    sources:
      - tor-browser
      - onion-sites
      - educational-content
      - tails

sources:
  tor-browser:
    type: rsync
    description: "Tor Browser (latest, all platforms)"
    version_detect:
      method: http-scrape
      url: "https://www.torproject.org/dist/torbrowser/"
      pattern: '>([0-9]+\.[0-9]+(\.[0-9]+)?)/'
      select: highest-semver
    rsync_base: "rsync://rsync.torproject.org/dist-mirror/torbrowser/{version}/"
    files:
      - "tor-browser-windows-x86_64-portable-{version}.exe"
      - "tor-browser-windows-x86_64-portable-{version}.exe.asc"
      - "tor-browser-macos-{version}.dmg"
      - "tor-browser-macos-{version}.dmg.asc"
      - "tor-browser-linux-x86_64-{version}.tar.xz"
      - "tor-browser-linux-x86_64-{version}.tar.xz.asc"
      - "tor-browser-android-aarch64-{version}.apk"
      - "tor-browser-android-aarch64-{version}.apk.asc"
    verify:
      method: gpg
      keyring: keys/torproject-signing.gpg
    stage_path: "software/tor-browser/"

  onion-sites:
    type: git
    description: "Real-world onion sites directory"
    git_url: "https://github.com/alecmuffett/real-world-onion-sites.git"
    shallow: true
    stage_path: "resources/onion-sites/"

  tails:
    type: rsync
    description: "Tails OS (latest stable)"
    rsync_base: "rsync.tails.net::amnesia-archive"
    verify:
      method: gpg
      keyring: keys/tails-signing.gpg
    stage_path: "software/tails/"

  educational-content:
    type: local
    description: "Bundled guides and documentation"
    local_path: "./content/guides/"
    stage_path: "guides/"

signing:
  enabled: true
  key_id: ""               # GPG key ID; prompted at build time if empty
  public_key: keys/oak-signing.pub
```

### Design Rationale

- **Declarative**: The entire build is config-driven. No code changes needed to add sources of existing types.
- **Version templating**: `{version}` placeholders in URLs/filenames resolve at runtime.
- **Tier budgets account for filesystem overhead**: Raw USB capacity minus FAT32/exFAT overhead.
- **Source/tier separation**: Sources defined once, referenced by name in tier lists.

---

## 4. Build Pipeline

`oak build --tier <tier>` runs the full quarterly update cycle:

```
1. LOAD CONFIG         Parse oak.yaml, resolve tier, compute release name
2. DETECT VERSIONS     Scrape latest Tor Browser version, Tails version, etc.
3. DOWNLOAD            rsync/git/http fetch into mirror/ (parallel, idempotent)
4. VERIFY UPSTREAM     Check GPG sigs (.asc/.sig) against bundled keyrings
5. STAGE               Copy selected files from mirror/ -> image/OAK-Q126/
                       Enforce tier size budget; abort if exceeded
6. ANNOTATE            Generate README.txt, README.html, MANIFEST.txt, VERSION.txt
                       Render Markdown guides -> standalone HTML
                       Build companion website (Go-native renderer)
7. PACKAGE             ZIP image/ -> dist/OAK-Q126-64GB.zip + .sha256
8. SIGN                GPG-sign the ZIP -> dist/OAK-Q126-64GB.zip.asc
```

Each step can run independently. `oak download` only fetches. `oak verify` only checks. This supports iterative development: download once, experiment with staging.

---

## 5. Image Layout (USB Drive Structure)

```
OAK-Q126/
├── README.txt                    # Plain text — first thing users see
├── README.html                   # Rich version with links to docs/
├── VERSION.txt                   # Release name, build date, source versions
├── MANIFEST.txt                  # SHA256 of every file (sha256sum -c compatible)
│
├── software/
│   ├── tor-browser/
│   │   ├── README.txt            # What is Tor Browser + install instructions
│   │   ├── windows/
│   │   │   ├── tor-browser-windows-x86_64-portable-X.Y.Z.exe
│   │   │   └── tor-browser-windows-x86_64-portable-X.Y.Z.exe.asc
│   │   ├── macos/
│   │   │   ├── tor-browser-macos-X.Y.Z.dmg
│   │   │   └── tor-browser-macos-X.Y.Z.dmg.asc
│   │   ├── linux/
│   │   │   ├── tor-browser-linux-x86_64-X.Y.Z.tar.xz
│   │   │   └── tor-browser-linux-x86_64-X.Y.Z.tar.xz.asc
│   │   └── android/
│   │       ├── tor-browser-android-aarch64-X.Y.Z.apk
│   │       └── tor-browser-android-aarch64-X.Y.Z.apk.asc
│   └── tails/                    # 32GB+ tiers only
│       ├── README.txt            # What is Tails + how to flash
│       ├── tails-amd64-X.Y.img
│       ├── tails-amd64-X.Y.img.sig
│       ├── tails-amd64-X.Y.iso
│       └── tails-amd64-X.Y.iso.sig
│
├── resources/
│   └── onion-sites/              # Real-world onion sites directory
│       ├── README.md
│       └── ...
│
├── guides/                       # Standalone HTML guides (no JS required)
│   ├── what-is-tor.html
│   ├── using-tor-browser.html
│   ├── what-is-tails.html
│   ├── privacy-basics.html
│   └── censorship-circumvention.html
│
├── docs/                         # Companion website (self-contained)
│   ├── index.html                # Start here
│   ├── about/
│   ├── getting-started/
│   ├── tools/
│   └── css/style.css
│
└── keys/
    ├── torproject-signing.gpg
    ├── tails-signing.gpg
    ├── oak-signing.pub           # OAK builder's public key
    └── README.txt                # How to verify signatures manually
```

### Key Decisions

- **README.txt at root** — universal, works on any OS with any text editor
- **README.html** — richer experience with links to the companion website
- **Platform subdirectories** under `software/tor-browser/` — non-technical users find their OS easily
- **`docs/` works from `file://`** — no JS routing, relative URLs only, no CDN dependencies
- **MANIFEST.txt** — offline-verifiable with `sha256sum -c MANIFEST.txt`
- **Separation**: `software/` (installable executables), `resources/` (reference data), `guides/` (educational reading)

---

## 6. Companion Website (Go-Native Rendering)

The `oak` CLI renders the companion website itself — no external tools required.

### Technology

- **goldmark** (Go Markdown parser) — Markdown to HTML
- **html/template** — page layouts
- Embedded CSS (single `style.css`, no external dependencies)

### How It Works

The `internal/site/` package:
1. Reads Markdown from `content/` and HTML templates from `content/templates/site/`
2. Renders a static site into `image/OAK-Q126/docs/`
3. All URLs are relative — works from `file://` protocol
4. No JavaScript required for navigation
5. Minimal, clean design readable in Tor Browser at default security settings

### Content

- **Homepage** — what OAK is, what's included, how to get started
- **About** — project background, philosophy
- **Getting Started** — step-by-step for first-time users
- **Tool pages** — dedicated page per included tool (Tor Browser, Tails)
- **Guides** — privacy basics, censorship circumvention, onion routing

---

## 7. Signing Strategy

### Upstream Verification

`oak verify` checks `.asc` and `.sig` files from Tor Project and Tails against their public keys bundled in `keys/`. The build fails if verification fails (overridable with `--force`).

### OAK Image Signing

After packaging, `oak build` GPG-signs the output ZIP:

1. Builder provides their GPG key ID via `oak.yaml` `signing.key_id` or `--sign-key` flag
2. Produces `dist/OAK-Q126-64GB.zip.asc` (detached signature)
3. Builder's public key is:
   - Embedded in the image at `keys/oak-signing.pub`
   - Published on the GitHub repository
4. **Online users** verify the sig against the GitHub-published key
5. **Offline recipients** verify against the embedded key (trust-on-first-use — if you trust the person who gave you the USB, you trust the key on it)

This is pragmatic for Q126. A more robust web-of-trust or keyserver model can follow.

---

## 8. GitHub README Structure

```
# 🌳 Open Access Kit (OAK)

One-line description

[Badges: release, license, CI]

## What is OAK?
2-3 paragraph explanation

## Quick Start
### Download Pre-Built Image
Links to latest release ZIPs per tier

### Build Your Own
go install + oak build example

## What's Inside
Table: content x tier matrix

## How to Use OAK
Brief instructions for USB recipients

## Building from Source
Prerequisites, build, configuration

## Release Schedule
Quarterly cadence, naming convention

## Contributing
How to add sources, contribute guides

## License
```

---

## 9. Q126 Release Scope

### In

| Area | Deliverables |
|------|-------------|
| CLI | `oak build`, `download`, `verify`, `stage`, `status`, `version` |
| Sources | Tor Browser (rsync + version detect), onion-sites (git), Tails (rsync), educational content (local) |
| Tiers | 16GB, 32GB, 64GB, max |
| Verification | GPG verification of upstream Tor Browser + Tails |
| Signing | GPG signing of output ZIP |
| Website | Go-native companion site renderer |
| Content | 3-5 educational guides (Tor, Tails, privacy, censorship) |
| Artifacts | MANIFEST.txt, VERSION.txt, README.txt/html |
| Repo | GitHub README, CI (GitHub Actions) |

### Out (Deferred)

- Torrent-based downloading
- Mullvad Browser, Signal APKs, other tools
- Internationalization / localized content
- Interactive TUI
- GitHub Pages hosting
- Automated quarterly release GitHub Action
- Delta updates between releases

### Size Estimates

| Content | Est. Size |
|---------|-----------|
| Tor Browser (all platforms + sigs) | ~1.2 GB |
| Onion sites directory | ~15 MB |
| Tails ISO + IMG + sigs | ~2.8 GB |
| Educational guides + website | ~25 MB |
| Keys, manifests, READMEs | ~1 MB |
| **16GB tier total** | **~1.3 GB** |
| **32GB+ tier total** | **~4.1 GB** |

Significant headroom in all tiers — intentional for future content additions.

---

## 10. Implementation Sequence

1. **Project scaffolding** — `go.mod`, cobra CLI skeleton, config loading, `oak version` + `oak status`
2. **Source interface + Tor Browser** — port `tor-mirror.sh` to Go (version detect, rsync, GPG verify)
3. **Remaining sources** — Git (onion-sites), rsync (Tails), local (educational content)
4. **Tier logic + staging** — size budgets, mirror->image copy, platform subdirs
5. **Annotation** — MANIFEST.txt, VERSION.txt, README generation, Markdown->HTML
6. **Site renderer** — Go-native companion website from Markdown + templates
7. **Packaging + signing** — ZIP creation, SHA256, GPG signing
8. **Full pipeline** — `oak build` orchestration, dry-run support
9. **Content + polish** — educational guides, GitHub README, CI

---

## 11. Verification Plan

| Test | Command |
|------|---------|
| Dry run | `oak build --tier 16 --dry-run` — correct source selection, no downloads |
| Download + verify | `oak download && oak verify` — upstream GPG checks pass |
| Full build | `oak build --tier 16` — inspect image layout matches spec |
| Offline docs | Open `image/OAK-Q126/docs/index.html` from `file://` in browser |
| Manifest check | `sha256sum -c image/OAK-Q126/MANIFEST.txt` |
| Signature check | `gpg --verify dist/OAK-Q126-16GB.zip.asc` |
| All tiers | Build 16/32/64/max, confirm size budgets respected |
