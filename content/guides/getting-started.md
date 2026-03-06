# Getting Started

→ [Back to Home](index.md) · [Resources](resources.md)

---

This guide is for someone who just received Open Access Kit removable media. Follow these
steps to start browsing safely and privately.

---

## Step 1: Install Tor Browser

Open the `software/tor-browser/` folder on this removable media and install the version
for your operating system:

| File | Operating System |
|------|-----------------|
| `tor-browser-windows-x86_64-portable-*.exe` | Windows |
| `tor-browser-macos-*.dmg` | macOS |
| `tor-browser-linux-x86_64-*.tar.xz` | Linux |
| `tor-browser-android-aarch64-*.apk` | Android |

**Windows / macOS:** Run the installer and follow the prompts. Tor Browser does not
require administrator access and installs only for your user account.

**Linux:** Extract the `.tar.xz` file, open the resulting folder, and run `./start-tor-browser.desktop` or the `Browser/start-tor-browser` script.

**Android:** You may need to allow installation from unknown sources in your device
settings. Go to **Settings → Security → Install unknown apps** and allow your file manager.

---

## Step 2: Connect to Tor

Launch Tor Browser. On the connection screen, click **Connect**.

If direct access to Tor is blocked in your country, click **Configure Connection** and
try one of the built-in **bridges** (Snowflake or obfs4 are good first choices).

Once connected, the Tor Browser address bar will show a purple onion icon.

---

## Step 3: Browse Safely

Tor Browser protects your anonymity at the network level, but a few habits will make
you significantly safer:

- **Always use Tor Browser** for sensitive browsing — do not open `.onion` links or
  sensitive sites in your regular browser.
- **Do not log in to personal accounts** (Google, Facebook, etc.) while using Tor if
  you need to remain anonymous — login ties your activity to your identity.
- **Keep Tor Browser updated.** New versions ship on this drive each quarter, or
  download from [torproject.org](https://www.torproject.org) when you have access.
- **Do not install browser extensions.** They can undermine Tor's anonymity protections.

### Finding onion sites

The [Onion Sites directory](resources/onion-sites/index.md) contains curated `.onion`
addresses for news outlets, privacy tools, and other services — all reachable through
Tor Browser.

---

## Step 4: Boot Tails (32 GB+ drives only)

Tails is a live operating system that runs entirely from this removable media and leaves
no trace on the computer you use. Use it when you need stronger guarantees than Tor Browser
alone can provide.

**To boot Tails:**

1. Restart the computer.
2. While the computer is starting up, press the boot menu key. This varies by
   manufacturer — common keys are **F12**, **F8**, **Esc**, or **Del**.
3. Select this removable media from the boot menu.
4. Tails will start and walk you through connecting to Tor.

> **Note:** Tails is only available on 32 GB and larger OAK distributions. The `software/tails/`
> folder on smaller distributions will be absent.

---

## Need Help?

- **Tor Browser support:** [tb-manual.torproject.org](https://tb-manual.torproject.org)
- **Tails documentation:** [tails.net/doc](https://tails.net/doc)
- **EFF Surveillance Self-Defense:** [ssd.eff.org](https://ssd.eff.org) — practical
  security advice for people of all threat levels.

→ Advanced: [Verify software signatures](verify.md)

→ [Back to Home](index.md) · [Resources](resources.md)
