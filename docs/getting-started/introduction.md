# Introduction

Envelope is a self-hosted, API-first, multi-tenant email platform. It gives every tenant ("account") on a shared deployment its own bearer token, its own hosted sending domains ("vhosts"), and its own quota, retention, and webhook configuration — while sharing one Postgres-backed installation underneath.

This documentation covers Envelope's **external interface only**: the REST management API, the webhook event contract, the SMTP and IMAP surfaces, and how to deploy and operate an instance. It does not describe Envelope's internal implementation — you don't need to know how it's built to integrate against it or run it.

## Who this is for

- **Integrators** — applications that need to send transactional or bulk email, receive inbound mail programmatically, or both, without depending on a third-party ESP.
- **Operators** — teams standing up and running their own Envelope deployment for internal use or as a platform for their own customers.

## The mental model

```
Account ("Acme Inc")
 ├─ Vhost (mail.acme.example)      ← owns a DKIM key, a quota, a retention window
 │   ├─ Mailbox (billing@mail.acme.example)   ← SMTP AUTH / IMAP login identity
 │   └─ Webhook subscription(s)
 └─ Vhost (notifications.acme.example)
     └─ Mailbox (alerts@notifications.acme.example)
```

An account is a business or tenant. It can own **any number** of vhosts — there's no one-to-one mapping. Every bearer token is scoped to an account, not to a single vhost, so one token can manage every vhost (and every mailbox, webhook, and further token) that account owns. See [Multi-Tenancy](../core-concepts/multi-tenancy.md) for the full ownership model.

## Two ways to send mail

- **REST** — `POST /accounts/:accountId/messages`, a single authenticated HTTP call with to/cc/bcc, HTML and text bodies, attachments, and custom headers. See [API Reference → Messages](../api-reference/messages.md).
- **SMTP submission** — authenticated `AUTH PLAIN` on port 587, the traditional path for existing mail-sending code (`smtplib`, `nodemailer`, an MTA relay, etc.). See [SMTP → Submission](../smtp/submission.md).

Both paths share the same DKIM signing, the same outbound queue, and fire the same `message.queued` webhook event — pick whichever fits your integration better, or use both.

## Receiving mail

Envelope receives mail over SMTP on its own MX-facing port and makes it available three ways: a `message.received` webhook fired the moment a message is accepted, standard IMAP4rev2 access to the mailbox, or a bulk per-vhost data export. There is no REST endpoint to list or fetch individual received messages — see [Known Limitations](../reference/known-limitations.md).

## Next steps

- [Quickstart](quickstart.md) — send your first email in three requests
- [Authentication](authentication.md) — the token model and self-serve onboarding sequence
- [Deploy Envelope](../deployment/deploy.md) — if you need to stand up your own instance first
