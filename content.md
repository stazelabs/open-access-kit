# OAK Educational Content Sources

Inventory of freely-licensed educational resources suitable for offline bundling in OAK. Each source has been evaluated for license compatibility, offline feasibility, and relevance to OAK's mission of providing privacy tools and knowledge on removable media.

---

## Source Inventory

| Source | License | Format | Fetch Method | OAK Tier | Est. Size |
|--------|---------|--------|-------------|----------|-----------|
| Tor Browser Manual | CC-BY-4.0 | Markdown (Lektor) | git clone | 16GB+ | ~20 MB |
| Tor Support Portal | CC-BY-4.0 | HTML (multilingual) | site mirror | 32GB+ | ~50 MB |
| Tor Training Materials | CC-BY-SA-4.0 | HTML/slides | site mirror | 32GB+ | ~30 MB |
| EFF Surveillance Self-Defense | CC-BY | HTML | site mirror | 16GB+ | ~40 MB |
| Security in a Box | CC-BY-SA-4.0 | Markdown | git clone | 16GB+ | ~25 MB |
| Digital First Aid Kit | CC-BY-4.0 | HTML (offline ZIP) | HTTP download | 16GB+ | ~15 MB |
| Privacy Guides | CC-BY-SA-4.0 | Markdown | git clone | 16GB+ | ~30 MB |
| Access Now Guides | CC-BY-3.0-US | HTML (GitLab) | git clone | 16GB+ | ~20 MB |
| Tails Documentation | CC-BY-SA-4.0 | HTML | site mirror | 32GB+ | ~40 MB |
| FLOSS Manuals: Bypassing Censorship | GPL/CC | PDF | HTTP download | 32GB+ | ~5 MB |
| Citizen Lab Circumvention Guide | Open access | PDF | HTTP download | 32GB+ | ~3 MB |
| OONI Educational Materials | CC-BY-SA-4.0 | HTML/data | site mirror | 64GB+ | ~100 MB |
| Access Now #KeepItOn Report | Public | PDF | HTTP download | 32GB+ | ~10 MB |
| Onion Services Ecosystem | CC-BY-3.0/MIT | HTML (GitLab) | git clone | 64GB+ | ~15 MB |
| Whonix Documentation | OSI-approved | Wiki/HTML | site mirror | max | ~80 MB |
| RSF Digital Security Resources | Open | HTML/PDF | site mirror | 64GB+ | ~20 MB |
| Real-World Onion Sites | MIT | Markdown | git clone | 16GB+ | ~15 MB |
| Wikipedia: Onion Routing | CC-BY-SA-3.0 | HTML | HTTP download | 32GB+ | ~2 MB |
| SaferJourno (Internews) | CC-BY-SA-4.0 | Markdown | git clone | 64GB+ | ~15 MB |

---

## Detailed Source Descriptions

### Tor Browser Manual
- **URL**: https://tb-manual.torproject.org/
- **Repository**: https://github.com/torproject/manual
- **License**: CC-BY-4.0
- **Content**: Installation, configuration, censorship circumvention settings, onion services, troubleshooting, mobile usage
- **Format**: Markdown source built with Lektor framework; produces static HTML
- **Offline feasibility**: Excellent -- clone repo, build with Lektor or render Markdown directly
- **Notes**: Fundamental companion to Tor Browser binaries. Should always ship alongside them.

### Tor Support Portal
- **URL**: https://support.torproject.org/
- **License**: CC-BY-4.0
- **Content**: FAQ-style docs covering Tor basics, browser features, troubleshooting, relay operation. Multilingual (Spanish, Farsi, German, Turkish, Russian, Ukrainian, Arabic, Chinese, Japanese)
- **Format**: HTML website
- **Offline feasibility**: Good -- can be mirrored. Multilingual support is valuable for OAK's international audience.

### Tor Training Materials
- **URL**: https://community.torproject.org/training/
- **License**: CC-BY-SA-4.0
- **Content**: Slides, guides, fanzines, and digital security guides for educators and community organizers
- **Format**: HTML, slides, PDF
- **Offline feasibility**: High -- open-licensed, designed for redistribution

### EFF Surveillance Self-Defense (SSD)
- **URL**: https://ssd.eff.org/
- **License**: CC-BY (explicitly permits printing and sharing)
- **Content**: Four sections -- Basics (surveillance countermeasures), Tool Guides (step-by-step), Further Learning (theory), Security Scenarios (role-based playlists for LGBTQ youth, journalists, activists, researchers)
- **Format**: HTML with printable versions
- **Offline feasibility**: Good -- explicitly permits offline sharing with attribution

### Security in a Box
- **URL**: https://securityinabox.org/en/
- **Repository**: https://github.com/securityinabox/siabguide
- **License**: CC-BY-SA-4.0
- **Content**: Modular guides on computer/phone security, encrypted storage, secure communication, privacy. Designed for human rights defenders, activists, journalists.
- **Format**: Markdown source on GitHub; HTML website
- **Offline feasibility**: Excellent -- explicitly designed for offline use and custom bundling ("remixed with just the parts that most suit your particular needs")
- **Maintained by**: Front Line Defenders + Tactical Technology Collective

### Digital First Aid Kit
- **URL**: https://digitalfirstaid.org/
- **Repository**: https://gitlab.com/metamorphosis-org-mk/dfak
- **License**: CC-BY-4.0
- **Content**: Self-diagnostic tools for people facing digital threats: malware, compromised accounts, phishing, DDoS. Emergency response focus.
- **Format**: HTML with official offline ZIP download (dfak-offline.zip)
- **Offline feasibility**: Excellent -- provides pre-built offline bundle

### Privacy Guides
- **URL**: https://www.privacyguides.org/
- **Repository**: https://code.privacyguides.dev/privacyguides/privacyguides.org
- **License**: CC-BY-SA-4.0
- **Content**: Community-driven privacy tool comparisons and recommendations. Covers Tor, VPNs, DNS, messaging, email. Has a Tor onion address.
- **Format**: Markdown source; static site
- **Offline feasibility**: Very high -- CC-licensed, GitHub-hosted, can be cloned and built

### Access Now Digital Security Guides
- **URL**: https://guides.accessnow.org/
- **License**: CC-BY-3.0-US
- **Content**: Tool guides for password management, secure communications, device security, censorship circumvention. Available in English, French, Spanish, Portuguese, Turkish, Tibetan, Burmese.
- **Format**: HTML; source on GitLab
- **Offline feasibility**: Good -- open source, clonable

### Tails Documentation
- **URL**: https://tails.net/doc/index.en.html
- **License**: CC-BY-SA-4.0 (graphics/assets); GPL-v3+ (source code)
- **Content**: Installation, usage, persistence, anonymity, security practices for Tails OS
- **Format**: HTML with multi-language support
- **Offline feasibility**: Moderate -- documentation can be mirrored; essential companion to Tails ISO

### FLOSS Manuals: How to Bypass Internet Censorship
- **URL**: https://archive.flossmanuals.net/bypassing-censorship/
- **License**: GPL / Creative Commons
- **Content**: Comprehensive manual covering proxies, VPNs, Tor, bridges, and advanced circumvention. Translated to Russian, Burmese, Arabic.
- **Format**: HTML (web archive), PDF
- **Offline feasibility**: Very high -- available as free PDF download

### Citizen Lab: Everyone's Guide to Bypassing Internet Censorship
- **URL**: https://citizenlab.ca/guides/everyones-guide-english.pdf
- **License**: Open access (publicly available PDF)
- **Content**: Practical circumvention guide from University of Toronto's Citizen Lab
- **Format**: PDF
- **Offline feasibility**: Very high -- direct PDF download

### OONI (Open Observatory of Network Interference)
- **URL**: https://ooni.org/
- **Data**: https://docs.ooni.org/data/
- **License**: CC-BY-SA-4.0 (educational materials); CC-BY-NC-SA-4.0 (data)
- **Content**: Country-level censorship reports, network interference documentation, research
- **Format**: HTML, API, data exports (S3 bucket)
- **Offline feasibility**: High -- country reports downloadable; data is CC-licensed
- **Notes**: NC clause on data license -- verify compliance for OAK's use case

### Access Now #KeepItOn Report
- **URL**: https://www.accessnow.org/campaign/keepiton/
- **Report**: https://www.accessnow.org/wp-content/uploads/2025/02/KeepItOn-2024-Internet-Shutdowns-Annual-Report.pdf
- **License**: Public (freely downloadable)
- **Content**: Annual report on internet shutdowns worldwide (283+ shutdowns in 39+ countries in 2024)
- **Format**: PDF
- **Offline feasibility**: Very high -- direct PDF download

### Onion Services Ecosystem
- **URL**: https://onionservices.torproject.org/
- **Repository**: https://gitlab.torproject.org/tpo/onion-services/
- **License**: CC-BY-3.0-US / MIT / GPL-v3 (varies by component)
- **Content**: Technical documentation for running onion services, applications, specifications
- **Format**: HTML, GitLab repos
- **Offline feasibility**: Moderate -- can clone repos; more technical/advanced content

### Whonix Documentation
- **URL**: https://www.whonix.org/wiki/Documentation
- **License**: OSI-approved
- **Content**: Tor configuration, bridge setup, onion services, anonymity practices for Whonix OS
- **Format**: Wiki-based HTML
- **Offline feasibility**: Medium-high -- can mirror wiki pages

### RSF Digital Security Resources
- **URL**: https://safety.rsf.org/
- **Content**: Digital security guides for journalists, censorship circumvention, Collateral Freedom program (restores access to 150+ censored news sites)
- **License**: Open educational resources
- **Format**: HTML, PDF
- **Offline feasibility**: Medium-high

### SaferJourno (Internews)
- **URL**: https://saferjourno.internews.org/
- **Repository**: https://github.com/OpenInternet/saferjourno
- **License**: CC-BY-SA-4.0
- **Content**: 6-module curriculum for media trainers: risk assessment, basic protection, mobile safety, data security, secure research, email protection
- **Format**: Markdown on GitHub
- **Offline feasibility**: Good -- designed for trainer distribution

### Wikipedia: Onion Routing
- **URL**: https://en.wikipedia.org/wiki/Onion_routing
- **Related**: https://en.wikipedia.org/wiki/List_of_Tor_onion_services
- **License**: CC-BY-SA-3.0
- **Format**: HTML
- **Offline feasibility**: High -- can export individual articles

---

## Summary URL List (for pipeline integration)

### Git Clone Sources

```
# Core docs (16GB+)
https://github.com/torproject/manual
https://github.com/securityinabox/siabguide
https://code.privacyguides.dev/privacyguides/privacyguides.org
https://github.com/alecmuffett/real-world-onion-sites.git

# Extended (32GB+)
https://gitlab.com/metamorphosis-org-mk/dfak

# Reference (64GB+)
https://gitlab.torproject.org/tpo/onion-services/ecosystem
https://github.com/OpenInternet/saferjourno
```

### HTTP Downloads (PDF/ZIP)

```
# Core (16GB+)
https://digitalfirstaid.org/dfak-offline.zip

# Extended (32GB+)
https://archive.flossmanuals.net/bypassing-censorship/bypassing-censorship.pdf
https://citizenlab.ca/guides/everyones-guide-english.pdf
https://www.accessnow.org/wp-content/uploads/2025/02/KeepItOn-2024-Internet-Shutdowns-Annual-Report.pdf
```

### Site Mirrors (wget/httrack)

```
# Core (16GB+)
https://ssd.eff.org/

# Extended (32GB+)
https://support.torproject.org/
https://community.torproject.org/training/
https://tails.net/doc/
https://guides.accessnow.org/

# Reference (64GB+)
https://ooni.org/reports/
https://safety.rsf.org/
https://onionservices.torproject.org/
```

### Wikipedia Articles (API export)

```
# Extended (32GB+)
https://en.wikipedia.org/wiki/Onion_routing
https://en.wikipedia.org/wiki/Tor_(network)
https://en.wikipedia.org/wiki/List_of_Tor_onion_services
https://en.wikipedia.org/wiki/Internet_censorship
https://en.wikipedia.org/wiki/Pluggable_transport
```

---

## Tier Allocation Summary

### 16GB Tier (~170 MB of educational content)
- Tor Browser Manual (repo)
- EFF Surveillance Self-Defense (mirror)
- Security in a Box (repo)
- Digital First Aid Kit (offline ZIP)
- Privacy Guides (repo)
- Access Now Guides (repo)
- Real-World Onion Sites (repo, already included)

### 32GB Tier (adds ~140 MB)
- Tor Support Portal (mirror, multilingual)
- Tor Training Materials (mirror)
- Tails Documentation (mirror)
- FLOSS Manuals: Bypassing Censorship (PDF)
- Citizen Lab circumvention guide (PDF)
- Access Now #KeepItOn report (PDF)
- Wikipedia articles (API export)

### 64GB Tier (adds ~150 MB)
- OONI reports and educational materials (mirror)
- Onion Services Ecosystem docs (repo)
- RSF digital security resources (mirror)
- SaferJourno curriculum (repo)

### Max Tier (adds ~80 MB)
- Whonix Documentation (mirror)
- Additional OONI data exports
- Extended Wikipedia coverage

---

## License Compliance Notes

All sources identified use permissive licenses that allow offline redistribution:

- **CC-BY-4.0 / CC-BY-3.0**: Redistribute with attribution. Must credit original authors.
- **CC-BY-SA-4.0 / CC-BY-SA-3.0**: Redistribute with attribution; derivative works must use same license.
- **CC-BY-NC-SA-4.0** (OONI data only): Non-commercial use. OAK is non-commercial, but verify compliance.
- **MIT**: Redistribute with license notice.
- **GPL-v3**: Redistribute with source availability.

OAK images must include an `ATTRIBUTION.txt` file crediting each bundled source with its name, URL, authors, and license. This file should be auto-generated during the `oak annotate` step.

---

## Next Steps

1. Validate all URLs are still live and content is current
2. Add these as source definitions in `oak.yaml` (new source types: `git`, `http-download`, `site-mirror`, `wikipedia-export`)
3. Estimate actual sizes by performing test downloads
4. Write `ATTRIBUTION.txt` template for license compliance
5. Prioritize which sources to implement first for Q126
