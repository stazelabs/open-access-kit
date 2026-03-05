# Open Access Kit (OAK)

A portable, offline-first collection of privacy tools and educational content for people facing censorship or surveillance.

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

## What is OAK?

Downloading Tor Browser or Tails OS requires internet access — exactly what may be restricted for people who need them most. OAK solves this by pre-packaging these tools onto USB thumbdrives that can be physically shared and replicated, no internet required.

OAK is a Go CLI that automates the quarterly update cycle: detect latest upstream versions, download everything, verify cryptographic signatures, assemble a tiered image, and package it for distribution.

## Quick Start

### Download a Pre-Built Image

Pre-built images are published on the [Releases](../../releases) page. Choose the tier that fits your USB drive:

| Tier | Drive Size | Contents |
|------|-----------|----------|
| 16GB | 16 GB | Tor Browser + guides + onion sites |
| 32GB | 32 GB | + Tails OS |
| 64GB | 64 GB | + Tails OS (more headroom) |
| max  | Any | Everything, no size limit |

### Build Your Own

```bash
go install github.com/open-access-kit/oak/cmd/oak@latest
oak build --tier 64
```

## What's Inside

| Content | 16GB | 32GB | 64GB | max |
|---------|------|------|------|-----|
| Tor Browser (all platforms) | Yes | Yes | Yes | Yes |
| Educational guides | Yes | Yes | Yes | Yes |
| Onion sites directory | Yes | Yes | Yes | Yes |
| Tails OS | - | Yes | Yes | Yes |

## How to Use OAK

If someone handed you an OAK USB drive:

1. Open `README.txt` at the root of the drive
2. Navigate to `software/tor-browser/` and find your operating system
3. Install Tor Browser and connect to the Tor network
4. Browse `guides/` for privacy and censorship circumvention guides
5. Visit `docs/index.html` for the full companion website (works offline)

## Building from Source

**Prerequisites**: Go 1.22+, rsync, gpg

```bash
git clone https://github.com/open-access-kit/oak
cd oak
make build

# Run the full pipeline
./oak build --tier 16 --dry-run      # preview what would happen
./oak build --tier 16                 # actually build
```

See [DESIGN.md](DESIGN.md) for full architecture documentation.

## Release Schedule

OAK follows a quarterly release cadence. Releases are named by quarter and year:

- Q126 = Q1 2026
- Q226 = Q2 2026
- etc.

Each release picks up the latest stable versions of all included tools.

## Contributing

- **Adding a source**: If it fits an existing source type (rsync, git, http, local), add a block to `oak.yaml` — no Go code required.
- **Educational content**: Add Markdown files to `content/guides/`. They are rendered to standalone HTML at build time.
- **New source types**: Implement the `source.Source` interface in `internal/source/`.

See [AGENTS.md](AGENTS.md) for guidance if you're using AI coding tools.

## License

Apache 2.0 — see [LICENSE](LICENSE).
