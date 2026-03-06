# OnionShare

→ [Back to Home](index.md) · [Resources](resources.md)

---

OnionShare lets you share files, receive files, host a website, and chat — all anonymously
over the Tor network. There is no central server: your computer becomes a temporary onion
service that only the people you invite can reach.

OnionShare is bundled on this removable media in the `software/onionshare/` folder
(32 GB and larger distributions only).

---

## Install OnionShare

Open `software/onionshare/` on this removable media and install the file for your
operating system:

| File | Operating System |
|------|-----------------|
| `OnionShare-win64-*.msi` | Windows |
| `OnionShare-*.dmg` | macOS |
| `OnionShare-*.flatpak` | Linux |

**Windows:** Run the `.msi` installer and follow the prompts. Administrator access is
required.

**macOS:** Open the `.dmg`, drag OnionShare to your Applications folder, and launch it.
macOS may ask you to approve the app from an unidentified developer — go to
**System Settings → Privacy & Security** and click **Open Anyway**.

**Linux (Flatpak):** Install with:
```
flatpak install --user OnionShare-*.flatpak
```
If Flatpak is not installed, your distribution's package manager can add it:
`sudo apt install flatpak` (Debian/Ubuntu) or `sudo dnf install flatpak` (Fedora).

---

## Verify the Installation (Optional)

Every file in `software/onionshare/` is accompanied by a `.asc` signature file.
You can verify the download against the OnionShare signing keys bundled in `keys/onionshare-signing.gpg`:

```
gpg --no-default-keyring --keyring /path/to/keys/onionshare-signing.gpg \
    --verify OnionShare-*.msi.asc OnionShare-*.msi
```

A good result shows `Good signature from "Micah Lee"`, `"Saptak Sengupta"`, or `"Miguel Jacq"`.
See [Verifying Signatures](verify.md) for a full walkthrough.

---

## Send Files

Use **Send** mode to share files with someone. OnionShare creates a temporary onion
address; anyone with that address can download your files through Tor Browser.

1. Launch OnionShare and select **Share Files**.
2. Drag files or folders into the window and click **Start sharing**.
3. Copy the `.onion` address shown and send it to your recipient through a secure channel
   (Signal, an encrypted email, etc.).
4. **Keep OnionShare open** until the recipient confirms they have downloaded everything.
   Closing it ends the transfer.

By default, OnionShare stops sharing after the first download. Uncheck **Stop sharing
after files have been sent** if multiple people need to download.

---

## Receive Files

Use **Receive** mode to let someone upload files to you without revealing your IP address
to them.

1. Select **Receive Files** in OnionShare.
2. Click **Start Receive Mode**.
3. Share the `.onion` address with the person sending you files.
4. Received files appear in your Downloads folder. You will see a notification when
   a transfer completes.

---

## Host a Website

Use **Host a Website** mode to serve static HTML files as an anonymous onion site.

1. Select **Host a Website**.
2. Add your HTML files (and any images, CSS, etc.).
3. Click **Start sharing**.
4. The `.onion` address is your website — share it with anyone who should be able to
   visit it through Tor Browser.

This is useful for publishing documents or resources when you need to protect both your
identity and the identity of your readers.

---

## Anonymous Chat

Use **Chat** mode for a temporary, encrypted group chat room that requires no account.

1. Select **Chat** in OnionShare.
2. Click **Start chat server**.
3. Share the `.onion` address with participants — they open it in Tor Browser.
4. Messages are end-to-end encrypted and leave no server logs. When you close OnionShare,
   the chat room ceases to exist.

---

## Security Tips

- **Share `.onion` addresses through a secure channel.** Anyone who intercepts the address
  can access your share. Use Signal, an encrypted email, or another private channel.
- **OnionShare must stay open.** Your computer is the server. If you close OnionShare or
  lose power, the session ends. For large transfers, plug in and disable sleep mode.
- **OnionShare protects your IP, not your identity.** Do not upload files containing
  personal metadata (EXIF data in photos, document author fields, etc.) if anonymity
  matters. Consider scrubbing metadata first with a tool like [MAT2](https://0xacab.org/jvoisin/mat2).
- **Use Tor Browser to access OnionShare links.** Your regular browser does not support
  `.onion` addresses.

---

## More Information

- OnionShare documentation: [docs.onionshare.org](https://docs.onionshare.org)
- OnionShare source code: [github.com/onionshare/onionshare](https://github.com/onionshare/onionshare)
- License: GPL-3.0

→ [Back to Home](index.md) · [Resources](resources.md)
