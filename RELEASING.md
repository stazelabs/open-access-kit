# Releasing OAK

This document covers how to build OAK, set up a GPG signing key, publish the companion website to Cloudflare Pages, and publish release archives to Cloudflare R2.

## Prerequisites

- Go 1.22 or later
- `gpg` (GnuPG) — for signing and verification
- `rsync` — for mirroring Tor Browser and Tails
- `git` — for cloning onion-sites and the Tor Browser Manual
- `wrangler` (Cloudflare CLI) — for Pages and R2 deployments

Install Wrangler:

```sh
npm install -g wrangler
wrangler login
```

## Building the CLI

```sh
git clone https://github.com/stazelabs/open-access-kit.git
cd open-access-kit
go build -o oak ./cmd/oak
```

Or install globally:

```sh
go install github.com/stazelabs/open-access-kit/cmd/oak@latest
```

Verify the build:

```sh
./oak version
```

## Setting Up a Signing Key

OAK signs release archives with GPG. You need a key pair before running a signed build.

### Generate a new key

If you do not already have a suitable key:

```sh
gpg --full-generate-key
```

Recommended settings: RSA 4096, no expiry (or a multi-year expiry), your release identity as the name and email.

### Find your key ID

```sh
gpg --list-secret-keys --keyid-format long
```

The key ID is the 16-character hex string after `rsa4096/`, for example `8F3DA1B2C4E5F6A7`.

### Export and commit the public key

```sh
gpg --armor --export YOUR_KEY_ID > keys/oak-signing.pub
```

Commit `keys/oak-signing.pub` to the repository. This key is bundled into every OAK image so offline recipients can verify archives without internet access.

### Configure oak.yaml

Set `signing.key_id` in `oak.yaml` to your key ID:

```yaml
signing:
  enabled: true
  key_id: "8F3DA1B2C4E5F6A7"
  public_key: keys/oak-signing.pub
```

Alternatively, pass `--sign-key` at build time to avoid committing the key ID.

## Managing Key Material with 1Password

1Password can store the GPG private key and its passphrase, keeping them out of dotfiles and off the filesystem between uses.

### Prerequisites

Install the 1Password CLI and desktop app, then enable CLI integration in **Settings → Developer → Integrate with 1Password CLI**:

```sh
brew install 1password-cli   # macOS
op --version
```

### Store the signing key in 1Password

Export your private key and store it as a 1Password document:

```sh
gpg --armor --export-secret-keys YOUR_KEY_ID > /tmp/oak-signing.asc
op document create /tmp/oak-signing.asc \
  --title "OAK Signing Key" \
  --vault "Private"
rm /tmp/oak-signing.asc
```

Store the passphrase as a password item:

```sh
op item create \
  --category login \
  --title "OAK Signing Key Passphrase" \
  --vault "Private" \
  --field "password=your-passphrase-here"
```

Store the reference paths (adjust vault/title to match):

```
op://Private/OAK Signing Key/
op://Private/OAK Signing Key Passphrase/password
```

### Import the key for a build session

Use an ephemeral keyring so the private key is never written to your default `~/.gnupg`:

```sh
export GNUPGHOME=$(mktemp -d)
trap "rm -rf $GNUPGHOME" EXIT

# Retrieve and import the private key
op document get "OAK Signing Key" | gpg --batch --import

# Retrieve the passphrase and pre-seed the agent
PASSPHRASE=$(op read "op://Private/OAK Signing Key Passphrase/password")
echo "$PASSPHRASE" | gpg --batch --yes --passphrase-fd 0 \
  --pinentry-mode loopback \
  --quick-set-expire YOUR_KEY_ID 0   # no-op refresh to load the key into agent
```

Then run the build as normal. The `EXIT` trap removes the temporary keyring automatically.

### Signing without an interactive passphrase prompt

Pass the passphrase via stdin using GPG's loopback pinentry mode so the build is non-interactive:

```sh
PASSPHRASE=$(op read "op://Private/OAK Signing Key Passphrase/password")
echo "$PASSPHRASE" | gpg --batch --yes \
  --passphrase-fd 0 \
  --pinentry-mode loopback \
  --armor --detach-sign dist/OAK-Q126-64GB.zip
```

`oak sign` invokes `gpg` directly, so you can pre-load the passphrase into the running agent before calling `oak build --sign` and GPG will pick it up without prompting.

### Complete sign-and-release script

```sh
#!/usr/bin/env bash
set -euo pipefail

RELEASE=Q126

# Ephemeral keyring
export GNUPGHOME=$(mktemp -d)
trap "rm -rf $GNUPGHOME" EXIT

# Load key and passphrase from 1Password
op document get "OAK Signing Key" | gpg --batch --import
PASSPHRASE=$(op read "op://Private/OAK Signing Key Passphrase/password")

# Pre-seed GPG agent so oak sign is non-interactive
echo "$PASSPHRASE" | gpg --batch --passphrase-fd 0 \
  --pinentry-mode loopback --list-secret-keys > /dev/null

# Build and sign all tiers
for tier in 16 32 64 max; do
  ./oak build --tier "$tier" --sign
done

echo "Build complete. Artifacts in dist/"
```

---

## Running a Build

### Full pipeline (all six steps)

```sh
./oak build --tier 64 --sign
```

This runs: download → verify → stage → annotate → package → sign.

Output artifacts land in `dist/`:

```
dist/OAK-Q126-64GB.zip
dist/OAK-Q126-64GB.zip.sha256
dist/OAK-Q126-64GB.zip.asc
```

### Common flags

| Flag | Purpose |
|------|---------|
| `--tier 16\|32\|64\|max` | Target tier (default: `64`) |
| `--sign` | GPG-sign the output ZIP |
| `--sign-key KEY_ID` | Override the signing key |
| `--skip-download` | Skip step 1 when the mirror is already current |
| `--dry-run` | Show what would happen without doing it |
| `--release Q126` | Override the release name (default: from `oak.yaml`) |

### Build all tiers

```sh
for tier in 16 32 64 max; do
  ./oak build --tier $tier --sign
done
```

### Incremental rebuilds

Once `mirror/` is populated you can skip re-downloading on subsequent runs:

```sh
./oak build --tier 64 --skip-download --sign
```

### Running individual steps

```sh
./oak download              # Step 1: populate mirror/
./oak verify                # Step 2: GPG-check mirrored content
./oak stage --tier 64       # Step 3: assemble image/
./oak annotate              # Step 4: generate VERSION/MANIFEST/README
./oak package               # Step 5: create dist/ ZIP + SHA256
./oak sign dist/OAK-Q126-64GB.zip   # Step 6: detached GPG signature
```

## Rendering the Companion Website

`oak site` renders `content/guides/` to `docs/` as static HTML:

```sh
./oak site
```

Preview locally by opening `docs/index.html` in a browser.

To point links at a custom base URL (useful for staged deployments):

```sh
./oak site --base-url https://your-project.pages.dev
```

## Publishing to Cloudflare Pages

The `docs/` directory is the static site deployed to Cloudflare Pages.

### First-time setup

1. Log in to the [Cloudflare dashboard](https://dash.cloudflare.com) and go to **Workers & Pages**.
2. Create a new Pages project connected to the `open-access-kit` GitHub repository.
3. Set the build configuration:
   - **Framework preset**: None
   - **Build command**: leave blank (site is pre-rendered; see below)
   - **Build output directory**: `docs`
4. Under **Environment variables**, set none (the site has no build-time dependencies).

Because `oak` must run to render the site, the recommended workflow is to render locally and push the output rather than using Cloudflare's build workers:

1. Render: `./oak site`
2. Commit `docs/` and push to the main branch.
3. Cloudflare Pages deploys automatically on push.

Alternatively, use `wrangler pages deploy` for a direct upload without a git push:

```sh
./oak site
wrangler pages deploy docs --project-name open-access-kit
```

### Custom domain

In the Pages project settings, add a custom domain under **Custom domains**. Cloudflare manages the DNS record automatically for domains on Cloudflare.

## Publishing Releases to Cloudflare R2

Release archives are stored in a Cloudflare R2 bucket, giving recipients a stable download URL independent of GitHub.

### Create the R2 bucket

```sh
wrangler r2 bucket create oak-releases
```

Or create it in the Cloudflare dashboard under **R2 Object Storage**.

Enable public access on the bucket so that download URLs work without authentication:
in the dashboard, go to the bucket settings and turn on **Public access**, or bind a custom domain to the bucket.

### Upload release artifacts

After a successful `oak build --sign`, upload the `dist/` artifacts:

```sh
RELEASE=Q126
TIER=64GB

wrangler r2 object put oak-releases/${RELEASE}/${TIER}/OAK-${RELEASE}-${TIER}.zip \
  --file dist/OAK-${RELEASE}-${TIER}.zip

wrangler r2 object put oak-releases/${RELEASE}/${TIER}/OAK-${RELEASE}-${TIER}.zip.sha256 \
  --file dist/OAK-${RELEASE}-${TIER}.zip.sha256

wrangler r2 object put oak-releases/${RELEASE}/${TIER}/OAK-${RELEASE}-${TIER}.zip.asc \
  --file dist/OAK-${RELEASE}-${TIER}.zip.asc
```

Upload all tiers with a shell loop:

```sh
RELEASE=Q126
for tier in 16GB 32GB 64GB Max; do
  for ext in .zip .zip.sha256 .zip.asc; do
    file="dist/OAK-${RELEASE}-${tier}${ext}"
    [ -f "$file" ] || continue
    wrangler r2 object put "oak-releases/${RELEASE}/${tier}/OAK-${RELEASE}-${tier}${ext}" \
      --file "$file"
  done
done
```

### Public download URLs

With public access enabled, objects are reachable at:

```
https://pub-<your-account-hash>.r2.dev/oak-releases/Q126/64GB/OAK-Q126-64GB.zip
```

If you bind a custom domain to the bucket (e.g., `releases.openaccess.tools`), the URL becomes:

```
https://releases.openaccess.tools/oak-releases/Q126/64GB/OAK-Q126-64GB.zip
```

Update the download links in the GitHub README and companion website to point at these URLs for each quarterly release.

## Quarterly Release Checklist

- [ ] Update `release:` in `oak.yaml` (e.g., `Q126` → `Q226`)
- [ ] Run `./oak build --tier all --sign` for all tiers
- [ ] Run `./oak site` to refresh the companion website
- [ ] Upload `dist/` artifacts to R2: `wrangler r2 object put ...`
- [ ] Deploy the updated site: `wrangler pages deploy docs --project-name open-access-kit`
- [ ] Tag the release commit: `git tag Q226 && git push --tags`
- [ ] Update download links in the README to point at the new R2 objects
