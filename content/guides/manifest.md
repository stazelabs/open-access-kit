# Release Manifest — Q1 2026 (Q126)

This page describes the contents of the Q126 distribution. For cryptographic checksums
of every file on the drive, see `MANIFEST.txt` in the root of the removable media and verify
with:

```
sha256sum -c MANIFEST.txt
```

---

## Software

| Software | Version | Platforms | Approx. Size | Description |
|----------|---------|-----------|--------------|-------------|
| **Tor Browser** | latest stable | Windows, macOS, Linux, Android | ~200 MB | The Tor Project's hardened browser. Routes all traffic through the Tor anonymity network. GPG-verified against the Tor Project signing key. |
| **Tails OS** | latest stable | x86-64 (bootable removable media) | ~1.7 GB | A live operating system that runs from the removable media and leaves no trace. Includes Tor Browser, Thunderbird, and additional privacy tools. GPG-verified against the Tails signing key. *(M and L tiers)* |
| **OnionShare** | latest stable | Windows, macOS, Linux | ~385 MB | Share files, host websites, and chat anonymously over the Tor network. GPG-verified against the OnionShare signing keys. *(M and L tiers)* |
| **Orbot** | latest stable | Android (arm64-v8a, armeabi-v7a) | ~78 MB | Routes all Android app traffic through Tor. GPG-verified against the Guardian Project signing keys. |

## Content

| Item | Source | Approx. Size | Description |
|------|--------|--------------|-------------|
| **Tor Browser Manual** | [torproject/manual](https://github.com/torproject/manual) | ~20 MB | Official Tor Browser documentation: installation, bridges, circumvention, onion services, troubleshooting. CC-BY-4.0. |
| **Security in a Box** | [securityinabox/siabguide](https://github.com/securityinabox/siabguide) | ~10 MB | Tool-by-tool digital security guidance from Frontline Defenders. CC-BY-SA-3.0. |
| **Privacy Guides** | [privacyguides/privacyguides.org](https://github.com/privacyguides/privacyguides.org) | ~50 MB | Community-maintained recommendations for privacy-respecting software. CC-BY-SA-4.0. |
| **Digital First Aid Kit** | [rarenet/dfak_2020](https://gitlab.com/rarenet/dfak_2020) | ~10 MB | Self-diagnostic tools for people facing digital threats. CC-BY-4.0. |
| **Tails Documentation** | [tails.net/doc](https://tails.net/doc/) | ~80 MB | Full Tails OS documentation — installation, usage, security. CC-BY-SA-4.0. *(M and L tiers)* |
| **SaferJourno** | [OpenInternet/saferjourno](https://github.com/OpenInternet/saferjourno) | ~30 MB | Digital security curriculum for journalists, by Internews. CC-BY-SA-4.0. *(M and L tiers)* |
| **Guides** | This repository | < 1 MB | Offline HTML documentation covering getting started, this manifest, a resource directory, and curated onion site listings. |
| **ZIM P0 — survival & medical** | [Kiwix / zimgit](zim-content.md) | ~1.0 GB | Emergency medicine, water, knots, food safety, disaster prep, children's encyclopedia, FreeCodeCamp, Simple Wiktionary. *(M and L tiers)* |
| **ZIM P1 — practical reference** | [Kiwix](zim-content.md) | ~2.5 GB | PhET simulations, Appropedia, Energypedia, Wikivoyage. *(M and L tiers)* |
| **ZIM P1 — Simple Wikipedia** | [Kiwix](zim-content.md) | ~3.2 GB | Simple English Wikipedia. *(M tier only — L uses full Wikipedia mini instead)* |
| **ZIM P2 — deep reference** | [Kiwix](zim-content.md) | ~22 GB | Full English Wikipedia mini, Wikibooks, WikiMed, Wikiversity, SuperUser. *(L tier only)* |
| **ZIM P3 — extended** | [Kiwix](zim-content.md) | 30 GB+ | TED-Ed, Gutenberg literature, Wikipedia nopic, Stack Exchange archives. *Deferred — future 64 GB+ tier.* |

---

## Tiers

OAK is distributed in multiple sizes to fit different drives:

| Tier | Size Budget | Includes Tails? | Extra Guides |
|------|-------------|-----------------|--------------|
| S | max 4 GB | No | — |
| M | max 16 GB | Yes | Tails Documentation, SaferJourno |
| L | max 32 GB | Yes | Tails Documentation, SaferJourno |

All tiers include Tor Browser (all platforms), Orbot, Tor Browser Manual, Security in a Box, Privacy Guides, Digital First Aid Kit, and these guides. Tails, OnionShare, Tails Documentation, SaferJourno, Kiwix, and ZIM P0/P1 content are included in M and L tiers. ZIM P2 (Wikipedia mini and additional reference) is L tier only.

---

## Verification

All software is verified against upstream GPG signing keys before packaging. The signing
keys for Tor Browser, Tails, OnionShare, and Orbot are bundled in `keys/` on this drive.
The OAK image itself is signed with the OAK release key.

→ [Step-by-step signature verification guide](verify.md)
