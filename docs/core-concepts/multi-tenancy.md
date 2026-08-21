# Multi-Tenancy

Envelope's data model has three levels, each owned by the one above it.

```
Account  ──1:N──▶  Vhost  ──1:N──▶  Mailbox
```

## Account

A business or tenant. Fields: `id`, `name`, `createdAt`. Created only via `POST /accounts` (admin-only, one-time per tenant). Owns bearer tokens and any number of vhosts — **there is no one-to-one mapping between an account and a vhost**. A single business with multiple sending domains (e.g. a transactional domain and a marketing domain) manages both under one account, with one set of tokens.

## Vhost

A hosted sending/receiving domain, e.g. `mail.acme.example`. Owned by exactly one account. Each vhost has:

- Its own DKIM key pair, generated automatically when the vhost is created
- Its own policy: `maxMessageBytes`, `dailyQuota`, `spamRejectThreshold`, `spamQuarantineThreshold`, `retentionDays`
- Its own webhook subscriptions
- An `active` flag — a deactivated vhost stops accepting new mail but its data and history remain queryable

There is no endpoint to delete a vhost outright — only [deactivate](../api-reference/vhosts.md) it (reversible-in-spirit, though there's no separate reactivate endpoint either) and separately [erase its stored message content](../api-reference/data-export.md) for compliance purposes. The vhost record, its mailboxes, and its account all persist.

## Mailbox

A send/receive identity on a vhost, e.g. `billing@mail.acme.example`. A mailbox is what SMTP `AUTH` and IMAP `LOGIN` authenticate against — its credential is a bcrypt-hashed password, separate from the account/vhost's bearer token system entirely. You only need a mailbox if you plan to use SMTP submission or IMAP; REST send and webhook receipt work without one.

A mailbox's credential can only ever send DKIM-signed mail as its own vhost's domain — there's no way to authenticate as one vhost and send `From:` a different one.

## Tokens are account-scoped, not vhost-scoped

A bearer token belongs to an account (`Token.accountId`) and can act on **every vhost that account owns**. This is deliberate: a business with five sending domains manages all five with one token (or a handful of tokens it mints itself, one per service or environment) rather than juggling a separate credential per domain.

<div class="callout gap">
There is currently no way to mint a token scoped to only one of an account's vhosts — every token is account-wide. If you need to hand a credential to a third party that should only touch one domain, put that domain under its own account instead.
</div>

## Ownership in practice

| Action | Requires |
|---|---|
| Create an account | Admin token |
| Create a vhost | Admin token, or that account's own token (self-serve) |
| Create a mailbox, subscribe a webhook, read the audit log, export/erase data | Admin token, or the token of the account that owns the vhost |
| Mint/revoke a token, send a message | Admin token, or that account's own token |

See [Authentication](../getting-started/authentication.md) for the full onboarding sequence and the exact authorization rules behind each of these.

## Next steps

- [DNS and DKIM](dns-and-dkim.md) — what publishing a vhost for real sending requires
- [Rate Limits and Quotas](rate-limits-and-quotas.md) — every limit that applies per account/vhost/IP
- [API Reference → Accounts](../api-reference/accounts.md) and [→ Vhosts](../api-reference/vhosts.md)
