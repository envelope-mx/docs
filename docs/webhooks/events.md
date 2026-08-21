# Events

## Envelope shape

Every delivery, regardless of event type, wraps its payload the same way:

```json
{
  "id": "evt_...",
  "type": "message.delivered",
  "vhost": "mail.acme.example",
  "createdAt": "2026-08-20T12:00:03Z",
  "data": { }
}
```

`id` is one UUID per logical event, shared across every subscription it fans out to — use it to deduplicate retried deliveries, since a retry resends the identical `id`. `vhost` is the domain string, not its internal ID.

## Event catalog

| Event | Fires | `data` shape |
|---|---|---|
| `message.received` | Once per distinct recipient vhost touched by an inbound SMTP transaction, on both `accept` and `quarantine` outcomes (never on `reject`) | `{ "from": "...", "to": ["..."], "quarantine": false }` |
| `message.queued` | Once per accepted outbound send request — SMTP submission or REST send, not once per recipient | `{ "from": "...", "to": ["...", "..."] }` |
| `message.delivered` | Once per successful delivery attempt to a remote MTA | `{ "from": "...", "to": "..." }` |
| `message.bounced` | Once per job that hit a permanent (5xx) SMTP failure, or exhausted its retry attempts | `{ "from": "...", "to": "...", "attempts": 4, "error": "550 5.1.1 user unknown" }` |
| `message.deferred` | Once per retry-eligible failed delivery attempt (temporary/4xx, or a transport error) that hasn't yet exhausted retries | `{ "from": "...", "to": "...", "attempts": 2, "error": "connection timed out", "nextAttemptAt": "..." }` |
| `message.complained` | **Never fires** — no ISP feedback-loop (FBL) ingestion exists. Defined so a subscription's `eventTypes` can reference it without erroring, not because it's implemented. | — |

<div class="callout note">
Note the asymmetry in <code>to</code>'s shape: <code>message.received</code> and <code>message.queued</code> carry it as an <strong>array</strong> (one inbound transaction or one send request can cover multiple recipients); <code>message.delivered</code>, <code>message.bounced</code>, and <code>message.deferred</code> carry it as a <strong>single string</strong>, since outbound delivery is tracked one recipient per job.
</div>

## What "queued" doesn't tell you

`message.queued` means Envelope accepted the send and staged it for delivery — it is not a delivery confirmation. Wait for `message.delivered` (success), `message.bounced` (permanent failure), or `message.deferred` (temporary failure, will retry) to know what actually happened at the destination.

## Next steps

- [Signatures](signatures.md) — verify a delivery is genuinely from Envelope before trusting its payload
- [Handle bounces and deferrals](../guides/handle-bounces-and-deferrals.md)
- [Receive and verify webhooks](../guides/receive-and-verify-webhooks.md)
