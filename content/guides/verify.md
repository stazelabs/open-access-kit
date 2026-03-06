# Verifying Software Signatures

This page is for advanced users who want to independently verify the authenticity of
software on this drive. Verification is optional — the OAK build pipeline checks
all signatures automatically before packaging — but it lets you confirm the software
has not been tampered with after it left the build system.

The GPG public keys you need are bundled in `keys/` on this drive.

---

## Verify Tor Browser

Tor Browser installers are signed by the Tor Project. Each installer has a corresponding
`.asc` signature file.

```
# Import the Tor Project signing key
gpg --import keys/torproject-signing.gpg

# Verify the installer (substitute your platform and version)
gpg --verify software/tor-browser/tor-browser-windows-x86_64-portable-*.exe.asc \
              software/tor-browser/tor-browser-windows-x86_64-portable-*.exe

gpg --verify software/tor-browser/tor-browser-macos-*.dmg.asc \
              software/tor-browser/tor-browser-macos-*.dmg

gpg --verify software/tor-browser/tor-browser-linux-x86_64-*.tar.xz.asc \
              software/tor-browser/tor-browser-linux-x86_64-*.tar.xz

gpg --verify software/tor-browser/tor-browser-android-aarch64-*.apk.asc \
              software/tor-browser/tor-browser-android-aarch64-*.apk
```

A successful verification ends with `Good signature from "Tor Browser Developers (signing key) <torbrowser@torproject.org>"`.

The canonical Tor Project signing key and its fingerprint are published at
[torproject.org/en/download/](https://www.torproject.org/en/download/) — cross-check if
you received this drive from an unknown source.

---

## Verify Tails

Tails images are signed by the Tails signing key. *(M and L tiers only.)*

```
# Import the Tails signing key
gpg --import keys/tails-signing.gpg

# Verify the image (substitute version)
gpg --verify software/tails/tails-amd64-*.img.sig \
              software/tails/tails-amd64-*.img

gpg --verify software/tails/tails-amd64-*.iso.sig \
              software/tails/tails-amd64-*.iso
```

The canonical Tails signing key and fingerprint are published at
[tails.net/install/download/](https://tails.net/install/download/).

---

## Verify the OAK image itself

If you received OAK as a ZIP archive rather than pre-loaded media, verify the archive
signature before extracting:

```
# Import the OAK signing key (bundled on the drive or available from GitHub)
gpg --import keys/oak-signing.pub

# Verify the archive
gpg --verify OAK-Q126-M.zip.asc OAK-Q126-M.zip
```

> **Note:** `keys/oak-signing.pub` is added to the image at release time. It will not
> be present in development builds or test images.

The OAK signing key is also published in the source repository at
[github.com/stazelabs/open-access-kit/tree/main/keys](https://github.com/stazelabs/open-access-kit/tree/main/keys).

---

## About the bundled keys

The keys in `keys/` are checked into the OAK source repository as trust anchors — they
are committed to git and auditable by anyone. Do not replace them without a corresponding
commit explaining the change.

If you want to verify the keys themselves, compare their fingerprints against the
official sources:

- **Tor Project key**: [torproject.org/en/download/](https://www.torproject.org/en/download/)
- **Tails key**: [tails.net/install/download/](https://tails.net/install/download/)
- **OAK key**: [github.com/stazelabs/open-access-kit](https://github.com/stazelabs/open-access-kit)

