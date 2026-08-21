# Glossary

**Account** — A business or tenant. Owns any number of vhosts and holds bearer tokens. Created only via `POST /accounts` (admin-only, one-time). See [Multi-Tenancy](../core-concepts/multi-tenancy.md).

**Admin token** — The platform bootstrap credential (`ENVELOPE_API_ADMIN_TOKEN`). Authorizes every account. Used for exactly two things in normal operation: creating a new account, and cross-tenant listing (`GET /accounts`, `GET /vhosts`). See [Authentication](../getting-started/authentication.md).

**Correlation ID** — An identifier shared across every log line touching one message's lifecycle, even as it crosses roles (inbound → filter → storage → webhook, or submission → queue → deliverer → webhook). Useful for tracing one message's path through operator logs.

**Dead-letter** — A webhook delivery that exhausted every retry attempt without a `2xx` response. Stops retrying automatically, but its full attempt history remains queryable and it can be manually [redriven](../api-reference/webhooks.md) at any time.

**Deliverer** — The role that drains the outbound queue, resolves MX records, and attempts real SMTP delivery to remote mail servers. Fires `message.delivered`, `message.bounced`, and `message.deferred`.

**DKIM selector** — The label prefixed to a domain's DKIM DNS record name (`<selector>._domainkey.<domain>`). Always the literal string `envelope` for every vhost — not configurable. See [DNS and DKIM](../core-concepts/dns-and-dkim.md).

**Mailbox** — A send/receive identity on a vhost (e.g. `billing@mail.acme.example`), authenticated with a password for SMTP `AUTH` and IMAP `LOGIN`. See [Multi-Tenancy](../core-concepts/multi-tenancy.md).

**Quarantine** — An inbound message the filter pipeline flagged as likely spam (or scored during an rspamd outage) but didn't outright reject. Delivered into a distinct `Quarantine` IMAP folder rather than `INBOX`; still fires `message.received`, with `quarantine: true`.

**Redrive** — Manually re-attempting a dead-lettered webhook delivery, using the subscription's current URL and secret. See [API Reference → Webhooks](../api-reference/webhooks.md).

**Token** — A bearer credential scoped to one account (`env_...`, shown once at creation, stored only as a SHA-256 hash). Can act on every vhost its account owns — there's no narrower per-vhost scope. See [Authentication](../getting-started/authentication.md).

**Vhost** ("virtual host") — A hosted sending/receiving domain, owned by exactly one account. Has its own DKIM key, quota, spam thresholds, retention window, and webhook subscriptions. See [Multi-Tenancy](../core-concepts/multi-tenancy.md).

## Next steps

- [Multi-Tenancy](../core-concepts/multi-tenancy.md)
- [Known Limitations](known-limitations.md)
