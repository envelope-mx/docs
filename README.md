# Envelope Docs

The source for Envelope's public documentation site — a small Go static-site generator (built on the [Goose framework](https://github.com/awesome-goose/goose), same shape as [awesome-goose/docs](https://github.com/awesome-goose/docs)) plus the Markdown content it builds.

Live site: **https://envelope-mx.github.io/docs/**

Requires Go 1.25+.

## Structure

```
main.go, app/       — the CLI itself (build/serve/help commands)
builder/             — page collection, HTML templates, CSS, JS (all plain Go — no bundler)
docs/                — the actual documentation content (Markdown, no front matter)
dist/                — generated output (gitignored)
```

Page titles come from each file's first `# Heading`. Section grouping/ordering and per-page ordering within a section are hardcoded in `builder/builder.go` (`organizeSections`, `pageOrderMaps`) — there's no nav config file to edit separately; add a page by dropping a `.md` file in the right `docs/<section>/` directory and, if it needs a specific position, adding it to that section's entry in `pageOrderMaps`.

## Development

```bash
go run . build              # writes dist/
go run . serve --port 3000  # serves dist/ at http://localhost:3000
```

Or build a binary once and reuse it:

```bash
go build -o envelope-docs .
./envelope-docs build
./envelope-docs serve --port 3000
```

## Deploying

`.github/workflows/deploy.yml` builds and publishes `dist/` to GitHub Pages on every `v*` tag push, using `github.event.repository.name` as the site's base URL — so it publishes to `https://envelope-mx.github.io/docs/` as long as this repository stays named `docs` under the `envelope-mx` org. It assumes this repository itself is public — if it isn't, either make it public or replace the workflow with your own deploy step (any static host works, since the output is plain files).

## Before publishing for real

One placeholder still needs a real value filled in:

- **Registry/release paths** in `docs/deployment/{binary,docker,docker-compose,kubernetes}.md` (`<registry>/envelope`, the GitHub releases URL) — Envelope's core product source isn't public, only its binaries and Docker images are, so these docs deliberately describe *consuming* a published artifact rather than building from source. Point them at wherever those are actually published.

GitHub links (sidebar footer, "Edit this page on GitHub") and the Go module path both already point at this repository's real remote, [envelope-mx/docs](https://github.com/envelope-mx/docs).

## Related repositories

| Repository | Description |
| --- | --- |
| [envelope-mx/envelope](https://github.com/envelope-mx/envelope) | Core platform (private — binaries and Docker images published, source is not public) |
| [envelope-mx/index](https://github.com/envelope-mx/index) | Product/technical planning docs this documentation is derived from |
| [envelope-mx/envelope-mx.github.io](https://github.com/envelope-mx/envelope-mx.github.io) | GitHub Pages entry point that redirects to this site |
