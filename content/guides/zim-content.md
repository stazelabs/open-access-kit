# ZIM File Priority Analysis for OAK

> **Reference appendix.** For a guide on using Kiwix and ZIM content, see
> [Kiwix & Offline Content](kiwix.md).

## Purpose

This document identifies and prioritizes Kiwix ZIM files for inclusion in OAK
releases. The goal is to maximize offline utility for people in disconnected or
crisis settings — those who need factual information, medical guidance, practical
survival knowledge, and educational resources — while respecting the tight space
budgets of removable media.

## Methodology

### Evaluation criteria

Each ZIM file was scored across four dimensions:

1. **Crisis relevance** — Does this content help someone survive, stay healthy,
   or navigate an emergency? Medical references, water purification, food
   safety, and disaster preparedness score highest.

2. **Knowledge breadth** — How much general-purpose factual information does
   this provide? Wikipedia variants and encyclopedic references score high.

3. **Practical skills** — Does the content teach actionable skills (agriculture,
   energy, computing, navigation)? Hands-on content outranks purely academic
   material.

4. **Size efficiency** — What is the utility-per-gigabyte? A 67 MB medical
   encyclopedia that could save lives scores infinitely higher per-byte than a
   168 GB video archive, however valuable the latter may be.

### Variant selection strategy

Kiwix publishes many ZIM files in multiple variants:

| Variant | Description | Typical size vs. maxi |
|---------|-------------|-----------------------|
| **maxi** | Full content with all images and media | 100% (baseline) |
| **nopic** | Full article text, no images | ~30-40% |
| **mini** | Reduced article set with minimal images | ~10-15% |
| **top** | Only highest-traffic articles, full media | varies |

**Our strategy**: prefer *nopic* for large reference works (Wikipedia,
Wiktionary) since text carries most informational value in crisis settings.
Use *maxi* for small, visually dependent content (medical atlases, PhET
simulations, children's content). Use *mini* when nopic isn't available or
when the full article set is too large for the target tier.

### Data sources

- Kiwix OPDS catalog: `https://library.kiwix.org/catalog/v2/entries`
- Download mirrors: `https://download.kiwix.org/zim/{category}/`
- File sizes verified against directory listings as of March 2026

---

## Key Tradeoffs

### Wikipedia: which variant?

The single most important decision. English Wikipedia ranges from 3.1 GB
(Simple English, maxi) to 115 GB (full, maxi):

| Variant | Size | Articles | Best for |
|---------|------|----------|----------|
| `wikipedia_en_simple_all_maxi` | 3.1G | ~240K | M tier, learners, non-native speakers |
| `wikipedia_en_top_maxi` | 7.6G | ~50K | Quick reference with images |
| `wikipedia_en_all_mini` | 11G | 6.9M+ | L tier, full coverage compressed |
| `wikipedia_en_all_nopic` | 48G | 6.9M+ | Future XL/Max tier, full text |
| `wikipedia_en_all_maxi` | 115G | 6.9M+ | Future Max tier only |

**Decision**: Simple English for M, mini for L. Larger variants deferred
until demand justifies a 64GB+ tier. Simple English is underrated — it's
written at an accessible reading level, making it ideal for non-native
speakers and young readers.

### Breadth vs. depth

Including one massive resource (Wikipedia nopic at 48G) vs. many smaller
ones is a constant tension. We favor **breadth at lower tiers** (many small,
high-impact ZIMs) and **depth at higher tiers** (larger comprehensive
references). An M-tier kit with medical guides, survival content, Simple
Wikipedia, and TED-Ed videos is more useful than one with just Wikipedia mini
and nothing else.

### Language coverage

This analysis focuses on English-language content. Multilingual and
non-English ZIMs are noted where they exist but are not included in tier
allocations. Future work should create language-specific profiles.

### Video content

Video ZIMs (TED, Khan Academy, Crash Course) are extremely large relative
to text. TED-Ed at 5.3G is the most space-efficient video offering. Khan
Academy (168G) and Crash Course (21G) are Max-tier only. Video is valuable
for education but text content wins on information density per byte.

---

## Space Budget

### Existing source usage (estimated)

| Source | Approx. size | Present in tiers |
|--------|-------------|------------------|
| Tor Browser (all platforms) | 400 MB | all |
| Tails (img + iso) | 3.0 GB | M, L |
| OnionShare (all platforms) | 200 MB | M, L |
| Orbot | 50 MB | all |
| Kiwix tools + desktop | 100 MB | M, L |
| Git repos (manuals, guides) | 500 MB | all |
| Keys, local content | 50 MB | all |

### Available ZIM budget per tier

| Tier | Total budget | Existing sources | Available for ZIM |
|------|-------------|-----------------|-------------------|
| S (4 GB hard cap) | 4 GB | ~1 GB | **~3 GB** |
| M (16 GB media) | 14 GB | ~4.5 GB | **~9.5 GB** |
| L (32 GB media) | 29 GB | ~4.5 GB | **~24.5 GB** |

Note: S tier has no Kiwix tools/desktop or ZIM files — it is a pure privacy
toolkit. ZIM content starts at the M tier, which includes kiwix-tools and
kiwix-desktop.

---

## Curated Priority List

Files are grouped by priority level. Within each level, they are ordered by
utility-per-byte (highest first).

### P0 — Critical (include in M and L tiers)

Tiny, high-impact survival and medical content. Total: ~1 GB.
S tier has no ZIM budget — P0 starts at M.

| # | File | Size | Category | Rationale |
|---|------|------|----------|-----------|
| 1 | `zimgit-medicine_en_2024-08.zim` | 67 MB | Medical | Offline medical encyclopedia; could save lives |
| 2 | `wikem_en_all_maxi_2021-02.zim` | 42 MB | Medical | Emergency medicine wiki (WikEM) — **not in OPDS catalog; removed from oak.yaml pending URL verification** |
| 3 | `zimgit-water_en_2024-08.zim` | 20 MB | Survival | Water purification and management |
| 4 | `zimgit-knots_en_2024-08.zim` | 27 MB | Survival | Practical knot-tying reference |
| 5 | `zimgit-food-preparation_en_2025-04.zim` | 93 MB | Survival | Food safety and preparation |
| 6 | `librepathology_en_all_maxi_2025-09.zim` | 76 MB | Medical | Pathology reference with images |
| 7 | `freecodecamp_en_all_2026-02.zim` | 7.6 MB | Education | Full programming curriculum, incredibly compact |
| 8 | `vikidia_en_all_maxi_2025-12.zim` | 66 MB | Education | Children's encyclopedia (age-appropriate) |
| 9 | `wiktionary_en_simple_all_nopic_2026-01.zim` | 25 MB | Reference | Simple English dictionary |
| 10 | `zimgit-post-disaster_en_2024-05.zim` | 615 MB | Survival | Disaster preparedness and recovery |

**P0 subtotal: ~1.04 GB**

### P1 — High priority (include starting at M tier)

Core reference and education content to fill the M-tier budget. Combined with
P0, targets ~9.5 GB.

| # | File | Size | Category | Rationale |
|---|------|------|----------|-----------|
| 11 | `wikipedia_en_simple_all_maxi_2025-11.zim` | 3.1 GB | Reference | Accessible English encyclopedia with images |
| 12 | `phet_en_all_2026-02.zim` | 102 MB | Education | Interactive science/math simulations |
| 13 | `appropedia_en_all_maxi_2026-02.zim` | 555 MB | Sustainability | Appropriate technology, permaculture, development |
| 14 | `energypedia_en_all_nopic_2025-12.zim` | 689 MB | Sustainability | Energy access and renewable energy |
| 15 | `wikivoyage_en_all_maxi_2025-12.zim` | 1.1 GB | Reference | World geography and travel knowledge |
| 16 | `ted_mul_ted-ed_2026-01.zim` | 5.3 GB | Education | Educational video lessons across subjects |
| 17 | `gutenberg_en_lcc-s_2025-12.zim` | 4.2 GB | Literature | Agriculture and rural life (practical) |

**P0 + P1 subtotal: ~12.2 GB** (exceeds M tier's ~9.5 GB ZIM budget — need
to drop TED-Ed or Gutenberg agriculture from M, or accept a tighter fit by
trimming P1 picks. See tier allocation summary below.)

### P2 — Medium priority (include starting at L tier)

Broader reference, structured learning, and community Q&A. Combined with
P0+P1, targets ~24.5 GB.

| # | File | Size | Category | Rationale |
|---|------|------|----------|-----------|
| 18 | `wikipedia_en_all_mini_2025-12.zim` | 11 GB | Reference | Full English Wikipedia, compressed (replaces Simple English) |
| 19 | `wikibooks_en_all_nopic_2026-01.zim` | 2.9 GB | Education | Textbooks on every subject |
| 20 | `mdwiki_en_all_maxi_2025-11.zim` | 2.1 GB | Medical | Comprehensive medical database with images |
| 21 | `wikiversity_en_all_maxi_2026-02.zim` | 2.2 GB | Education | University-level learning materials |
| 22 | `superuser.com_en_all_2026-02.zim` | 3.7 GB | Tech | Computer troubleshooting Q&A |
| 23 | `gutenberg_en_lcc-h_2025-12.zim` | 4.2 GB | Literature | Social sciences (economics, sociology) |

**P0 + P1 + P2 subtotal: ~23.5 GB** (adjustments below)

**L tier note**: At L, `wikipedia_en_all_mini` (11G) replaces
`wikipedia_en_simple_all_maxi` (3.1G) from P1, so the net add is ~7.9G,
not 11G. Adjusted total: ~20.4 GB within 24.5 GB budget. Remaining ~4 GB
could accommodate additional P2+ picks from the expansion list below.

### P3 — Extended (deferred — future 64 GB+ tier)

Comprehensive references and richer educational content. Not included in
current S/M/L tiers. Revisit when demand justifies a larger tier.

| # | File | Size | Category | Rationale |
|---|------|------|----------|-----------|
| 24 | `wikipedia_en_all_nopic_2025-12.zim` | 48 GB | Reference | Complete English Wikipedia, all articles (replaces mini) |
| 25 | `wiktionary_en_all_nopic_2026-02.zim` | 8.2 GB | Reference | Complete English dictionary |
| 26 | `math.stackexchange.com_en_all_2026-02.zim` | 6.9 GB | Education | Mathematics Q&A archive |
| 27 | `electronics.stackexchange.com_en_all_2026-02.zim` | 3.9 GB | Tech | Electronics engineering Q&A |
| 28 | `ted_mul_science_2026-02.zim` | 13 GB | Education | Science-focused TED talks |
| 29 | `ted_mul_health_2026-01.zim` | 7.0 GB | Medical | Health-focused TED talks |
| 30 | `gutenberg_en_lcc-q_2025-12.zim` | 16 GB | Literature | Science books |
| 31 | `gutenberg_en_lcc-t_2025-12.zim` | 12 GB | Literature | Technology books |
| 32 | `wikibooks_en_all_maxi_2026-01.zim` | 5.1 GB | Education | Textbooks with images (replaces nopic) |

**Future 64 GB tier note**: If published, `wikipedia_en_all_nopic` (48G)
would replace `wikipedia_en_all_mini` (11G), net add ~37G.
`wikibooks_en_all_maxi` (5.1G) replaces nopic (2.9G), net add ~2.2G. A
64 GB tier couldn't fit everything in P3 — prioritize Wikipedia nopic +
Wiktionary + one TED collection + one Gutenberg collection. Estimated fit:
items 24-27 + item 29 ≈ ~52 GB ZIM total.

### P4 — Future Max tier only

Large archives that provide tremendous depth but require hundreds of
gigabytes. Deferred until demand justifies publishing a Max tier.

| # | File | Size | Category | Rationale |
|---|------|------|----------|-----------|
| 33 | `wikipedia_en_all_maxi_2026-02.zim` | 115 GB | Reference | Complete Wikipedia with all images |
| 34 | `stackoverflow.com_en_all_2023-11.zim` | 75 GB | Tech | Complete Stack Overflow archive |
| 35 | `ted_mul_all_2025-08.zim` | 79 GB | Education | All TED talks, all languages |
| 36 | `crashcourse_en_all_2026-02.zim` | 21 GB | Education | Crash Course video series |
| 37 | `gutenberg_en_all_2025-11.zim` | 206 GB | Literature | Complete Project Gutenberg (70K+ books) |
| 38 | `khanacademy_en_all_2023-03.zim` | 168 GB | Education | Complete Khan Academy |
| 39 | `wikipedia_es_all_maxi_2026-02.zim` | 38 GB | Reference | Spanish Wikipedia |
| 40 | `wikipedia_ar_all_maxi_2026-02.zim` | 17 GB | Reference | Arabic Wikipedia |
| 41 | `wikipedia_fr_all_maxi_2026-02.zim` | ~50 GB | Reference | French Wikipedia |
| 42 | `wikipedia_de_all_maxi_2026-01.zim` | 49 GB | Reference | German Wikipedia |

---

## Expansion List (ranks 43-100)

Additional ZIM files for consideration, especially at higher tiers or for
language-specific builds. Ordered by estimated priority.

### Medical & health (43-48)

| # | File | Size | Rationale |
|---|------|------|-----------|
| 43 | `africanstorybook.org_mul_all_2025-01.zim` | 8.1 GB | Multilingual children's literacy |
| 44 | `skinofcolorsociety_en_all_maxi_2024-05.zim` | 821 MB | Dermatology reference, diverse skin tones |
| 45 | `who.int_en_all_2025-02.zim` | ~500 MB | WHO health guidelines (if available) |
| 46 | `mdwiki_en_all_2025-11.zim` | 10 GB | Full MDWiki (larger variant with more media) |

### Practical skills & sustainability (47-55)

| # | File | Size | Rationale |
|---|------|------|-----------|
| 47 | `gutenberg_en_lcc-s_2025-12.zim` | 4.2 GB | Agriculture (already in P1) |
| 48 | `gutenberg_en_lcc-g_2025-12.zim` | 7.5 GB | Geography and anthropology |
| 49 | `gutenberg_en_lcc-e_2026-03.zim` | 9.5 GB | American history |
| 50 | `gutenberg_en_lcc-d_2025-12.zim` | 37 GB | World history |
| 51 | `gutenberg_en_lcc-pz_2025-12.zim` | 18 GB | Fiction and literature |
| 52 | `gutenberg_en_lcc-pr_2025-12.zim` | 15 GB | English literature |
| 53 | `gutenberg_en_lcc-n_2025-12.zim` | 21 GB | Fine arts |
| 54 | `gutenberg_en_lcc-a_2026-03.zim` | 9.1 GB | General reference works |
| 55 | `gutenberg_en_lcc-m_2025-12.zim` | 3.7 GB | Music |

### Education & video (56-65)

| # | File | Size | Rationale |
|---|------|------|-----------|
| 56 | `ted_mul_social-change_2026-01.zim` | 10 GB | Social change talks |
| 57 | `ted_mul_technology_2026-01.zim` | 15 GB | Technology talks |
| 58 | `ted_mul_education_2026-01.zim` | 7.2 GB | Education talks |
| 59 | `ted_mul_psychology_2026-01.zim` | 4.8 GB | Psychology talks |
| 60 | `ted_mul_ted-conference_2026-02.zim` | 17 GB | Main TED conference talks |
| 61 | `ted_mul_business_2026-01.zim` | 8.9 GB | Business and economics talks |
| 62 | `ted_mul_sustainability_2026-02.zim` | 5.4 GB | Sustainability talks |
| 63 | `ted_mul_history_2026-01.zim` | 5.4 GB | History talks |
| 64 | `ted_mul_culture_2026-01.zim` | 10 GB | Culture talks |
| 65 | `ted_mul_creativity_2026-01.zim` | 5.1 GB | Creativity talks |

### Stack Exchange communities (66-75)

| # | File | Size | Rationale |
|---|------|------|-----------|
| 66 | `tex.stackexchange.com_en_all_2026-02.zim` | 4.2 GB | LaTeX/document preparation |
| 67 | `blender.stackexchange.com_en_all_2026-02.zim` | 2.6 GB | 3D modeling |
| 68 | `askubuntu.com_en_all_2026-02.zim` | ~5 GB | Ubuntu/Linux help |
| 69 | `serverfault.com_en_all_2026-02.zim` | ~4 GB | Server administration |
| 70 | `dba.stackexchange.com_en_all_2026-02.zim` | ~2 GB | Database administration |
| 71 | `unix.stackexchange.com_en_all_2026-02.zim` | ~4 GB | Unix/Linux Q&A |
| 72 | `physics.stackexchange.com_en_all_2026-02.zim` | ~3 GB | Physics Q&A |
| 73 | `chemistry.stackexchange.com_en_all_2026-02.zim` | ~2 GB | Chemistry Q&A |
| 74 | `biology.stackexchange.com_en_all_2026-02.zim` | ~2 GB | Biology Q&A |
| 75 | `security.stackexchange.com_en_all_2026-02.zim` | ~3 GB | Information security (relevant to OAK mission) |

### Reference & languages (76-90)

| # | File | Size | Rationale |
|---|------|------|-----------|
| 76 | `wiktionary_fr_all_nopic_2026-02.zim` | ~3 GB | French dictionary |
| 77 | `wiktionary_es_all_nopic_2026-02.zim` | ~2 GB | Spanish dictionary |
| 78 | `wikivoyage_de_all_maxi_2026-01.zim` | 1.2 GB | German travel guide |
| 79 | `wikivoyage_fr_all_maxi_2025-12.zim` | 273 MB | French travel guide |
| 80 | `wikipedia_en_simple_all_nopic_2025-11.zim` | ~1.5 GB | Simple English nopic (lighter alternative) |
| 81 | `wikiversity_de_all_maxi_2026-01.zim` | 1.3 GB | German educational content |
| 82 | `wikiversity_fr_all_maxi_2026-02.zim` | 488 MB | French educational content |
| 83 | `gutenberg_fr_all_2026-01.zim` | 9.8 GB | French literature |
| 84 | `gutenberg_de_all_2026-01.zim` | 10 GB | German literature |
| 85 | `gutenberg_es_all_2026-01.zim` | 1.7 GB | Spanish literature |
| 86 | `wikibooks_de_all_maxi_2026-01.zim` | 3.4 GB | German textbooks |
| 87 | `wikibooks_fr_all_maxi_2026-01.zim` | 1.9 GB | French textbooks |
| 88 | `wikibooks_es_all_maxi_2025-10.zim` | 333 MB | Spanish textbooks |
| 89 | `wikipedia_pt_all_maxi_2026-02.zim` | ~20 GB | Portuguese Wikipedia |
| 90 | `wikipedia_ru_all_maxi_2026-02.zim` | ~25 GB | Russian Wikipedia |

### Specialized & niche (91-100)

| # | File | Size | Rationale |
|---|------|------|-----------|
| 91 | `chopin.lib.uchicago.edu_en_all_2025-01.zim` | 8.2 GB | Classical music archive |
| 92 | `freecodecamp_es_all_2026-02.zim` | 7.6 MB | Spanish programming |
| 93 | `freecodecamp_fr_all_2026-02.zim` | 7.6 MB | French programming |
| 94 | `freecodecamp_pt_all_2026-02.zim` | 7.6 MB | Portuguese programming |
| 95 | `phet_mul_all_2026-02.zim` | 227 MB | PhET simulations (all languages) |
| 96 | `vikidia_fr_all_maxi_2025-12.zim` | 966 MB | French children's encyclopedia |
| 97 | `vikidia_es_all_maxi_2025-12.zim` | 91 MB | Spanish children's encyclopedia |
| 98 | `khanacademy_es_all_2023-03.zim` | 150 GB | Spanish Khan Academy |
| 99 | `gutenberg_mul_all_2025-11.zim` | 236 GB | Multilingual complete Gutenberg |
| 100 | `wikipedia_zh_all_maxi_2026-02.zim` | ~22 GB | Chinese Wikipedia |

---

## Tier Allocation Summary

### S tier (4 GB hard cap — no ZIM)

The S tier is a pure privacy toolkit. No Kiwix tools or ZIM files are
included. Total content is ~1 GB (Tor Browser, Orbot, guides, keys),
leaving comfortable margin within the 4 GB hard cap.

### M tier (~9.5 GB ZIM budget)

| Priority | Files | ZIM total |
|----------|-------|-----------|
| P0 | Items 1-10 | 1.04 GB |
| P1 | Items 11-14 (Simple Wikipedia, PhET, Appropedia, Energypedia) | 4.4 GB |
| P1 | Item 15 (Wikivoyage) | 1.1 GB |
| **Total** | **15 ZIM files** | **~6.5 GB** |

Margin: ~3 GB. Can selectively add TED-Ed (5.3 GB, item 16) if it fits, or
Gutenberg agriculture (4.2 GB, item 17) — but not both. Alternatively, keep
the margin for quarterly growth.

### L tier (~24.5 GB ZIM budget)

| Priority | Files | ZIM total |
|----------|-------|-----------|
| P0 | Items 1-10 | 1.04 GB |
| P1 | Items 12-17 (keep, except Simple Wiki replaced by mini) | 7.9 GB |
| P2 | Items 18-23 | 15.1 GB |
| **Total** | **22 ZIM files** | **~20.4 GB** |

Margin: ~4.1 GB. Can add 1-2 items from expansion list (e.g., items 44+55
= 4.5 GB, or item 59 = 4.8 GB).

**L tier note**: At this tier, `wikipedia_en_all_mini` (11G) replaces
`wikipedia_en_simple_all_maxi` (3.1G) from P1, so the net add is ~7.9G.

### Future tiers (64 GB+, Max)

Not currently published. If demand warrants it, a 64 GB tier could add
Wikipedia nopic (48G) and a Max tier could include the full P3/P4 lists.
See P3 and P4 sections above for candidate content.

---

## Integration Notes

### Relationship to oak.yaml

ZIM files would be added as new sources in `oak.yaml`, likely as `type: http`
sources downloading from `download.kiwix.org`. Each ZIM file (or group) would
be a separate source entry with appropriate `stage_path` under a `zim/`
directory.

Kiwix tools and desktop are included in the M and L tiers. The S tier has
no ZIM content and therefore does not need Kiwix software.

### Suggested stage layout

```
zim/
  reference/
    wikipedia_en_*.zim
    wiktionary_en_*.zim
    wikivoyage_en_*.zim
  medical/
    zimgit-medicine_en_*.zim
    wikem_en_*.zim
    mdwiki_en_*.zim
    librepathology_en_*.zim
  survival/
    zimgit-post-disaster_en_*.zim
    zimgit-water_en_*.zim
    zimgit-knots_en_*.zim
    zimgit-food-preparation_en_*.zim
  education/
    freecodecamp_en_*.zim
    phet_en_*.zim
    ted_mul_*.zim
    wikibooks_en_*.zim
    wikiversity_en_*.zim
    vikidia_en_*.zim
  sustainability/
    appropedia_en_*.zim
    energypedia_en_*.zim
  literature/
    gutenberg_en_*.zim
  tech/
    superuser.com_en_*.zim
    stackoverflow.com_en_*.zim
    *.stackexchange.com_en_*.zim
```

### Version management

ZIM filenames include dates (e.g., `_2026-02`). The download process should
fetch the latest available version. This could be implemented via:
- OPDS catalog query to resolve current filename
- Directory listing scrape of `download.kiwix.org/zim/{category}/`
- Pattern matching similar to existing `version_detect` in oak.yaml

### Verification

Kiwix provides SHA-256 checksums alongside ZIM files. The download source
type should verify checksums after download.

---

## Open Questions

1. **Language profiles**: Should OAK support language-variant builds (e.g.,
   an Arabic-focused L tier)? This would change the ZIM selection
   significantly.

2. **Update frequency**: ZIM files are updated monthly to quarterly. Should
   OAK pin specific versions or always fetch latest?

3. **Demand for larger tiers**: If users request 64 GB+ builds, which
   Wikipedia variant should anchor it — nopic (48G, all articles, no images)
   or mini (11G) plus many more supplementary resources?
