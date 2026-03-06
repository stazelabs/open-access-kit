# Contributing to Open Access Kit

Thank you for helping make privacy tools more accessible. Contributions of all kinds are
welcome — new sources, guides, bug fixes, and documentation improvements.

## Ways to Contribute

### Add a new source

Sources are defined in [`oak.yaml`](oak.yaml). OAK supports four source types: `rsync`,
`git`, `http`, and `local`. To add a new software or content source:

1. Fork the repository and create a branch.
2. Add an entry under `sources:` in `oak.yaml` with the appropriate type, URL, and
   verification method.
3. Add the source name to the relevant `tiers:` entries.
4. Run `oak download <source-name>` and `oak verify <source-name>` to confirm it works.
5. Open a pull request describing the source, its license, and why it belongs in OAK.

All bundled content must be freely redistributable. See the license notes in
[`content/guides/resources.md`](content/guides/resources.md) for examples of compatible licenses.

### Add or improve a guide

Guides live in [`content/guides/`](content/guides/) as Markdown files. They are rendered
to HTML and bundled on every distribution of removable media, as well as published to the
project website.

1. Add or edit a `.md` file in `content/guides/`.
2. Use relative links between guides (e.g., `[Getting Started](getting-started.md)`).
3. Run `./oak site` to preview the rendered output in `docs/guides/`.
4. Open a pull request with a brief description of the new or updated content.

### Fix a bug or improve the CLI

1. Read [`ARCHITECTURE.md`](ARCHITECTURE.md) to understand the build pipeline.
2. Run `go test ./...` to verify the existing tests pass.
3. Make your change, add tests if appropriate, and open a pull request.

## Code Style

- Standard Go formatting (`gofmt`). The CI linter (`golangci-lint`) enforces this.
- Keep changes minimal and focused. One concern per PR.
- Configuration-driven where possible — prefer extending `oak.yaml` over hard-coding.

## License

By contributing, you agree that your contributions will be licensed under the same terms
as the project: [GPL v3](LICENSE) for code, [CC BY-SA 4.0](LICENSE) for content.
