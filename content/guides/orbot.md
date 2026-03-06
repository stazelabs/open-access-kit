# Orbot

Orbot routes your Android device's internet traffic through the Tor network. Where Tor
Browser for Android anonymizes only your browsing, Orbot can protect other apps too —
messaging clients, email, podcast apps, and more — by tunneling their connections through
Tor as well.

Orbot is included in the `software/orbot/` folder on this removable media.

---

## Which APK to Install

Open `software/orbot/` on this removable media. You will find two APK variants:

| File | For |
|------|-----|
| `Orbot-*-arm64-v8a-release.apk` | Most Android devices (2014 and newer, 64-bit ARM) |
| `Orbot-*-armeabi-v7a-release.apk` | Older or budget Android devices (32-bit ARM) |

If you are unsure which to use, install the `arm64-v8a` version. If it reports
"App not installed" during setup, try the `armeabi-v7a` version instead.

---

## Install Orbot

1. On your Android device, go to **Settings → Security** (or **Privacy**) and enable
   **Install from unknown sources** (or **Install unknown apps** for the file manager you'll use).

2. Copy the APK file from this removable media to your device (via USB cable, SD card, or
   a local file transfer app).

3. Tap the APK file in your device's file manager and confirm the installation.

4. Disable "Install from unknown sources" after installation to reduce risk from other apps.

---

## Enable Orbot

Launch Orbot and tap **Start**. Orbot will connect to the Tor network.

### VPN Mode (Recommended)

Enable **VPN Mode** to route all apps through Tor automatically, without configuring each
app individually. Orbot will appear as a VPN in your Android status bar.

> **Note:** VPN mode routes all traffic through Tor. Apps that try to detect and block Tor
> (such as some banking or streaming apps) may not work while Orbot's VPN is active.

### Per-App Selection

Alternatively, use **Choose Apps** to select only specific apps to route through Tor,
leaving other apps on your regular internet connection.

---

## Orbot and Tor Browser for Android

**Tor Browser for Android does not need Orbot.** It includes its own Tor connection and
is already fully anonymized. Running both simultaneously is harmless but redundant.

Use Orbot for apps that do not have built-in Tor support — messaging apps, email clients,
news readers, and so on.

---

## Verify the APK (Optional)

Each APK in `software/orbot/` is accompanied by a `.asc` signature file. You can verify
it against the Guardian Project signing keys bundled in `keys/orbot-signing.gpg`:

```
gpg --no-default-keyring --keyring /path/to/keys/orbot-signing.gpg \
    --verify Orbot-*-arm64-v8a-release.apk.asc \
             Orbot-*-arm64-v8a-release.apk
```

A good result shows `Good signature from "Hans-Christoph Steiner"` or
`Good signature from "Nathan of Guardian"`.

See [Verifying Signatures](verify.md) for a full walkthrough.

---

## Security Tips

- **Orbot protects network-level identity, not app-level identity.** An app can still
  reveal your identity through your account, device fingerprint, or content of messages.
- **Keep Orbot running** while using apps you want protected. When Orbot is stopped or
  disconnected, those apps will fall back to your regular internet connection.
- **In high-risk situations, use Tails instead.** Tails (included on M and L tier distributions)
  routes all traffic system-wide through Tor by default, with stronger isolation guarantees.
- **Keep Orbot updated.** New versions ship on this removable media each quarter.

---

## More Information

- Orbot website: [orbot.app](https://orbot.app)
- Source code: [github.com/guardianproject/orbot-android](https://github.com/guardianproject/orbot-android)
- Guardian Project: [guardianproject.info](https://guardianproject.info)
- License: GPL-3.0
