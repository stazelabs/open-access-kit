# AGENTS.md — Guide for AI Coding Agents

This document gives AI coding agents the context needed to work effectively on the Open Access Kit (OAK) codebase. Read ARCHITECTURE.md for the full architecture specification.

## Project Summary

OAK is a Go CLI (`oak`) that curates, downloads, verifies, stages, and packages quarterly releases of privacy tools (Tor Browser, Tails OS) and educational content onto portable storage devices. The entire build pipeline is driven by `oak build`.

Target users: people in censorship/surveillance environments who need offline copies of privacy tools.

## Repository Layout

```
open-access-kit/
├── ARCHITECTURE.md         # Full architecture spec — read this first
├── AGENTS.md               # This file
├── oak.yaml                # Declarative config: tiers, sources, signing
├── go.mod / go.sum
├── Makefile
│
├── cmd/oak/main.go         # CLI entrypoint (cobra)
├── cli/                    # Cobra command definitions (build, download, verify, stage, status, version)
│
├── internal/
│   ├── config/             # oak.yaml loading and structs
│   ├── source/             # Source interface + rsync/git/http/local implementations
│   ├── version/            # Version detection (HTTP scrape, RSS)
│   ├── verify/             # GPG + checksum verification
│   ├── tier/               # Size budgets, content selection
│   ├── stage/              # Mirror -> image layout
│   ├── annotate/           # README/MANIFEST/VERSION generation
│   ├── packaging/          # ZIP creation + GPG signing
│   ├── site/               # Markdown->HTML renderer (goldmark + html/template)
│   └── pipeline/           # Full build orchestration
│
├── content/guides/         # Educational Markdown source files
├── content/templates/      # README and site HTML templates
├── keys/                   # Upstream GPG public keys (torproject, tails)
└── .github/workflows/      # CI (build, test, lint)
```

Runtime artifacts (gitignored): `mirror/`, `image/`, `dist/`

## Key Commands

```bash
make build          # go build ./cmd/oak
make test           # go test ./...
make lint           # golangci-lint run
make clean          # remove oak binary, mirror/, image/, dist/

# CLI usage
oak build --tier 64              # full pipeline
oak download [source]            # fetch/mirror sources
oak verify [source]              # GPG + checksum verification
oak stage --tier 64              # build image from mirror
oak status                       # show mirror state and versions
oak version                      # print CLI version
```

## Code Conventions

- **Go idioms**: standard Go error handling (`if err != nil`), no panics in library code
- **Packages**: each `internal/` subdirectory is one package; avoid cross-package circular deps
- **Config-driven**: adding a new source of an existing type should require only `oak.yaml` changes, not new Go code
- **Context propagation**: all I/O operations accept `context.Context` as first arg for cancellation
- **Interfaces**: `source.Source` interface is the core abstraction — implementations in `source/rsync.go`, `source/git.go`, etc.
- **No global state**: pass config/deps explicitly
- **Tests**: table-driven tests preferred; use `t.TempDir()` for filesystem operations

## Security-Sensitive Areas

The following areas require extra care:

- **`internal/verify/`** — GPG signature verification. Never skip or weaken verification logic. Failures must be hard errors (not warnings) by default.
- **`internal/packaging/`** — ZIP creation and GPG signing. Ensure no path traversal in archive entries.
- **`keys/`** — Public keys for Tor Project and Tails. Do not modify these files; they are the trust anchors.
- **Version detection** (`internal/version/`) — Scrapes upstream URLs. Be defensive about parsing; validate that detected versions are sane semver strings before using them in rsync paths.
- **`oak.yaml` signing config** — `signing.key_id` should never be hardcoded; it is prompted at build time.

When writing or reviewing code in these areas, be explicit about what could go wrong and add comments explaining the security invariants.

## Architecture Decisions (from ARCHITECTURE.md)

- **Cobra** for CLI framework
- **goldmark** for Markdown-to-HTML (companion website rendered natively — no external tools)
- **rsync protocol** for Tor Browser and Tails mirroring (idempotent, resumable)
- **Declarative config**: tiers reference source names; sources define their own type and verification
- **`{version}` templating** in rsync URLs/filenames resolved at runtime after version detection
- **Offline-first docs**: all HTML uses relative URLs, no JS routing, no CDN — works from `file://`

## What NOT to Do

- Do not add dependencies casually — this is a tool that should be easy to build from source
- Do not skip GPG verification or add `--force` behavior without explicit user request
- Do not store private keys or signing credentials in the repo
- Do not add JavaScript to the companion website — it must work at Tor Browser's "Safest" security level
- Do not add platform-specific build constraints unless truly necessary

## Testing Approach

```bash
# Dry run (no downloads)
oak build --tier 16 --dry-run

# Unit tests
go test ./internal/...

# Integration test (requires network + rsync)
oak download tor-browser && oak verify tor-browser
```

For unit tests of source implementations, use interfaces and mock the network layer. The `source.Source` interface makes this straightforward.
