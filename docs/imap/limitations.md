# IMAP Limitations

Envelope's IMAP surface is intentionally minimal — read access to what inbound filtering delivered, not a full mail-client backend. Know these before pointing a general-purpose IMAP client at it.

## Explicitly unsupported commands

`CREATE`, `DELETE` (mailbox), `RENAME`, `SUBSCRIBE`/`UNSUBSCRIBE`, `APPEND`, `IDLE`, `EXPUNGE`, `SEARCH`, `COPY`, `MOVE` — each returns IMAP `NO`, not a silent no-op.

<div class="callout gap">
No <code>APPEND</code> means there's no way to inject a message via IMAP — messages only ever arrive through <a href="../smtp/inbound.html">SMTP inbound</a> or the <a href="../smtp/submission.html">submission → queue → deliverer</a> path. No <code>IDLE</code> means no push notifications — a client has to poll (<code>SELECT</code>/<code>STATUS</code> again) to notice new mail; there's no live update feed.
</div>

## Session model

`SELECT` snapshots the mailbox at that moment — sequence numbers and UIDs are assigned by snapshot order and are **only stable for the lifetime of that one selection**. There's no `CONDSTORE`/`QRESYNC` persistence across sessions. If your client caches UIDs across reconnects expecting them to stay stable, verify it re-syncs on every new session rather than trusting a cached mapping.

## No search, no server-side filtering

`SEARCH` isn't implemented — a client has to `FETCH` and filter client-side if it needs anything beyond "everything in this folder." For most integrations, a [`message.received` webhook](../webhooks/events.md) is a better fit than polling IMAP for new mail in the first place; treat IMAP as a fallback/manual-inspection surface rather than the primary integration path.

## Next steps

- [Overview](overview.md)
- [Receive and verify webhooks](../guides/receive-and-verify-webhooks.md) — the recommended primary path for reacting to inbound mail
