# Known Limitations

Honest accounting of what Envelope's interface doesn't yet do, gathered in one place instead of only as scattered inline notes. None of these are secret — each is also called out on its relevant reference page.

## Authorization and accounts

- **No per-vhost token scope.** Every token is scoped to its whole account — there's no way to hand out a credential limited to only one of an account's vhosts. If you need that isolation, put the vhost under its own account instead.
- **No self-serve account creation.** Creating the first account for a new tenant is always an admin-only, one-time operation.
- **No account or vhost deletion.** Only [`PATCH /vhosts/:id/deactivate`](../api-reference/vhosts.md) (stop accepting mail) and [`DELETE /vhosts/:id/data`](../api-reference/data-export.md) (erase stored message content) exist — the records themselves are permanent.
- **No mailbox password rotation.** Delete and recreate the mailbox to change its credential.
- **No token update.** A token can only be created or revoked — there's no way to change its label after minting.

## Messages and sending

- **No message list/fetch API.** There is no REST endpoint to list or read individual received messages — use a [`message.received` webhook](../webhooks/events.md), [IMAP](../imap/overview.md), or the bulk [data export](../api-reference/data-export.md).
- **No message-size limit on REST send.** [`POST /accounts/:accountId/messages`](../api-reference/messages.md) enforces the daily sending quota but not `maxMessageBytes`, unlike SMTP submission.
- **No idempotency keys.** A retried `POST` after a dropped response is not automatically deduplicated on any endpoint — build your own idempotency at the caller if that matters to you.
- **Uniqueness conflicts surface as `500`, not `409`.** Creating a vhost for a domain that already exists, or a mailbox local-part that already exists on a vhost, currently returns a generic `500 Internal Server Error` rather than a REST-idiomatic `409 Conflict`.

## Webhooks

- **`message.complained` never fires.** It's defined so a subscription's `eventTypes` can reference it without erroring, but there's no ISP feedback-loop (FBL) ingestion behind it.
- **No subscription update.** A webhook subscription's URL and event types can't be changed after creation — disable it and create a new one instead.

## IMAP

- **No `APPEND`, `IDLE`, `SEARCH`, `EXPUNGE`, `COPY`, `MOVE`, or folder management.** See [IMAP Limitations](../imap/limitations.md) for the full list and what to use instead of each.

## DNS

- **No DNS readiness check.** Nothing confirms your DKIM/SPF/MX records have actually propagated — verify with an external DNS lookup tool.
- **No SPF or DMARC generation.** Envelope only manages DKIM; SPF and DMARC records are entirely your responsibility to write and publish.

## Deployment

- **No documented downgrade path.** Rolling back to a previous release against an already-migrated database isn't something this platform guarantees works — keep a database backup before upgrading production.

None of the above should block a first integration — they're worth knowing so you don't design around a capability that doesn't exist yet, not reasons to avoid the platform.

## Next steps

- [Errors and Responses](../core-concepts/errors-and-responses.md)
- [API Reference Overview](../api-reference/overview.md)
