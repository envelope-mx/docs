# Envelope

Envelope is a self-hosted, API-first, multi-tenant email platform. It sends and receives mail through a full-featured REST API or authenticated SMTP, signs everything with DKIM automatically, and reports on delivery through signed, retried webhooks — running entirely on infrastructure you control.

Envelope is built around three concepts:

- **Account** — a business or tenant. Owns any number of vhosts. Holds bearer tokens.
- **Vhost** — a hosted sending/receiving domain (e.g. `mail.acme.example`), owned by exactly one account. Has its own DKIM key, quota, and retention policy.
- **Mailbox** — a send/receive identity on a vhost (e.g. `alice@mail.acme.example`), used for SMTP AUTH and IMAP login.

One admin bootstrap token creates the first account for a new tenant. From there, everything — vhosts, mailboxes, further tokens, webhooks — is self-serve using that account's own token, with no further admin involvement. See [Multi-Tenancy](core-concepts/multi-tenancy.md) for the full model.

## Where to start

- New to Envelope? Start with [Introduction](getting-started/introduction.md), then follow the [Quickstart](getting-started/quickstart.md) to send your first email in a few requests.
- Standing up your own instance? See [Deploy Envelope](deployment/deploy.md).
- Integrating against a running instance? Jump straight to the [API Reference](api-reference/overview.md), [Webhooks](webhooks/overview.md), [SMTP](smtp/inbound.md), or [IMAP](imap/overview.md) reference.
