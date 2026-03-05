This directory contains GPG public keys used to verify upstream software.

Keys included:
- torproject-signing.gpg  — Tor Project release signing key
- tails-signing.gpg       — Tails OS release signing key
- oak-signing.pub         — OAK builder's public key (added at build time)

These keys are checked into the repository as trust anchors.
Do NOT modify them without careful review and a corresponding commit explaining why.

To manually verify a file:
  gpg --keyring ./keys/torproject-signing.gpg --verify tor-browser-*.asc
  gpg --keyring ./keys/tails-signing.gpg --verify tails-*.sig
