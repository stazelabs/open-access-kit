# 🌳 Open Access Kit (OAK)

A portable, offline-first collection of privacy tools and educational content for people
living under censorship or surveillance. Everything you need fits on removable media.

**Latest Release: Q1 2026 (Q126)** — [Download from GitHub](https://github.com/stazelabs/open-access-kit/releases/latest)

[What's Included](manifest.md) · [Getting Started](getting-started.md) · [Resources](resources.md)

---

## For End Users

OAK distributions are sized for the removable media you have available. All tiers include
Tor Browser for every platform and these guides.

| | 16 GB | 32 GB | 64 GB | Max |
|---|:---:|:---:|:---:|:---:|
| Tor Browser (Win / macOS / Linux / Android) | ✓ | ✓ | ✓ | ✓ |
| Tor Browser Manual (offline HTML) | ✓ | ✓ | ✓ | ✓ |
| Security in a Box (offline HTML) | ✓ | ✓ | ✓ | ✓ |
| Privacy Guides (offline HTML) | ✓ | ✓ | ✓ | ✓ |
| Digital First Aid Kit (offline HTML) | ✓ | ✓ | ✓ | ✓ |
| Orbot (Android) | ✓ | ✓ | ✓ | ✓ |
| Guides (this folder) | ✓ | ✓ | ✓ | ✓ |
| Tails OS (bootable live system) | — | ✓ | ✓ | ✓ |
| OnionShare (Win / macOS / Linux) | — | ✓ | ✓ | ✓ |
| Tails Documentation (offline HTML) | — | ✓ | ✓ | ✓ |
| SaferJourno (offline HTML) | — | — | ✓ | ✓ |

→ [Full manifest with versions and checksums](manifest.md)

### Getting started

1. Open the `software/tor-browser/` folder on this removable media.
2. Install Tor Browser for your operating system.
3. Launch Tor Browser and click **Connect**.
4. You're connected. Browse `.onion` sites using the
   [Onion Sites directory](resources/onion-sites/index.md), or read the guides in this folder.

→ [Full step-by-step instructions](getting-started.md)

---

## For Builders & Contributors

OAK is a Go CLI tool that downloads, verifies, stages, and packages quarterly offline
distributions of privacy tools. It is entirely configuration-driven via `oak.yaml`, and
designed to be auditable, reproducible, and easy to extend with new sources or guides.

- [Architecture](architecture.md) — how the build pipeline works
- [Contributing](contributing.md) — how to add sources, guides, or fixes
- [GitHub](https://github.com/stazelabs/open-access-kit) — source code, releases, and issue tracker
