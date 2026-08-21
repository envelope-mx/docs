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

GitHub links (sidebar footer, "Edit this page on GitHub"), the Go module path, and every registry/release reference in `docs/deployment/{binary,docker,docker-compose,kubernetes}.md` all point at the real published locations: [envelope-mx/docs](https://github.com/envelope-mx/docs) for this repository's own remote, [ghcr.io/envelope-mx/envelope](https://github.com/envelope-mx/envelope/pkgs/container/envelope) for the Docker image, and [envelope-mx/envelope releases](https://github.com/envelope-mx/envelope/releases) for binaries.

## Related repositories

| Repository | Description |
| --- | --- |
| [envelope-mx/envelope](https://github.com/envelope-mx/envelope) | Core platform (private — binaries and Docker images published, source is not public) |
| [envelope-mx/index](https://github.com/envelope-mx/index) | Product/technical planning docs this documentation is derived from |
| [envelope-mx/envelope-mx.github.io](https://github.com/envelope-mx/envelope-mx.github.io) | GitHub Pages entry point that redirects to this site |
