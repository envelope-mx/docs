# Envelope Docs

The source for Envelope's public documentation site — a small Go static-site generator (built on the [Goose framework](https://github.com/awesome-goose/goose), same shape as [awesome-goose/docs](https://github.com/awesome-goose/docs)) plus the Markdown content it builds.

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

`.github/workflows/deploy.yml` builds and publishes `dist/` to GitHub Pages on every `v*` tag push. It assumes this repository itself is public — if it isn't, either make it public or replace the workflow with your own deploy step (any static host works, since the output is plain files).

## Before publishing for real

A few placeholders need real values filled in:

- **Registry/release paths** in `docs/deployment/{binary,docker,docker-compose,kubernetes}.md` (`<registry>/envelope`, the GitHub releases URL) — Envelope's core product source isn't public, only its binaries and Docker images are, so these docs deliberately describe *consuming* a published artifact rather than building from source. Point them at wherever those are actually published.
- **GitHub links** in `builder/templates.go` (sidebar footer, "Edit this page on GitHub") currently point at `github.com/isaiahiroko/envelope-docs` — confirm that's this repository's real remote before publishing, or swap it for a "Suggest an edit" contact link if this repo won't be public either.
- **Module path** in `go.mod` (`github.com/isaiahiroko/envelope-docs`) if this ends up living somewhere else.
