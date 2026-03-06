# Kiwix & Offline Content

Kiwix is an offline content reader that lets you browse Wikipedia, educational
resources, medical references, and more — without any internet connection. OAK
bundles Kiwix along with curated ZIM files on M and L tier distributions.

> Kiwix and ZIM content are only available on **M and L** tier distributions.
> If you have the S tier, this section does not apply.

---

## What are ZIM files?

ZIM is an open format for storing web content offline. Each `.zim` file is a
self-contained archive of a website — Wikipedia, a medical encyclopedia, coding
tutorials, and so on. Kiwix reads these files and lets you search and browse
them just like you would online.

---

## Opening Kiwix

### Windows

1. Open the `software/kiwix-desktop/` folder on this removable media.
2. Extract `kiwix-desktop_windows_x64_*.zip`.
3. Run `kiwix-desktop.exe` from the extracted folder.

### Linux

1. Open the `software/kiwix-desktop/` folder on this removable media.
2. Make the AppImage executable: `chmod +x kiwix-desktop_x86_64_*.appimage`
3. Run: `./kiwix-desktop_x86_64_*.appimage`

### macOS / Linux (command line)

Kiwix-serve provides a browser-based reader:

1. Open the `software/kiwix/` folder on this removable media.
2. Extract the archive for your platform.
3. Start the server pointing at the ZIM directory:

```
./kiwix-serve --port 8080 /path/to/zim/**/*.zim
```

4. Open `http://localhost:8080` in any browser.

---

## Browsing ZIM content

Once Kiwix is open:

1. Click **Open file** (or use File > Open).
2. Navigate to the `zim/` folder on this removable media.
3. Select any `.zim` file to open it.
4. Use the search bar to find articles within the loaded content.

You can open multiple ZIM files and switch between them using the tabs or
library view.

---

## What's included

| Priority | Category | Examples | Tier |
|----------|----------|----------|------|
| **P0** — Survival & medical | Emergency medicine, water, food safety, disaster prep, children's encyclopedia, FreeCodeCamp | zimgit-medicine, librepathology, vikidia, freecodecamp | M, L |
| **P1** — Reference & education | Simple English Wikipedia, PhET simulations, Appropedia, Wikivoyage, TED-Ed | wikipedia\_en\_simple, phet, ted\_ed | M, L |
| **P2** — Deep reference | Full Wikipedia mini, Wikibooks, MDWiki medical, Wikiversity, SuperUser | wikipedia\_en\_all\_mini, mdwiki, wikibooks | L only |

ZIM files are organized in subdirectories under `zim/` by topic: `medical/`,
`survival/`, `education/`, `reference/`, `sustainability/`, `literature/`, and
`tech/`.

> For the full catalog of ZIM files with sizes and selection criteria, see the
> [ZIM Content Catalog](zim-content.md).

---

## Tips

- **Search is fast.** Kiwix indexes every article. Use the search bar to jump
  directly to topics.
- **Kiwix-serve for shared access.** Run `kiwix-serve` on one computer and let
  others on the same local network browse the content from their own browsers.
- **Mobile.** Kiwix is also available for Android and iOS. You can copy
  individual `.zim` files to a phone and open them in the Kiwix app.

> [Full manifest with versions and sizes](manifest.md)
