# Guide: Handle Bounces and Deferrals

After a message is queued (`message.queued`), the deliverer attempts real delivery to the recipient's MX and reports back via one of three events. Understanding the difference determines whether you should give up, wait, or investigate.

## The three outcomes

| Event | Meaning | What to do |
|---|---|---|
| `message.delivered` | Accepted by the destination MTA | Nothing — this recipient is done |
| `message.deferred` | A temporary (4xx) failure or transport error; will retry automatically | Usually nothing — wait for the next attempt (`data.nextAttemptAt`). Surface it in a dashboard if you're tracking delivery health, but don't treat it as final |
| `message.bounced` | A permanent (5xx) failure, or every retry attempt was exhausted | Final — stop retrying yourself too, and act on `data.error` |

```json
// message.deferred
{ "from": "billing@mail.acme.example", "to": "customer@example.com", "attempts": 2, "error": "421 4.3.0 try again later", "nextAttemptAt": "2026-08-20T12:15:00Z" }

// message.bounced
{ "from": "billing@mail.acme.example", "to": "customer@example.com", "attempts": 8, "error": "550 5.1.1 user unknown" }
```

## Interpreting `data.error`

The `error` field is the raw response text from the destination MTA (or a transport-level error like a connection timeout) — there's no structured bounce-reason taxonomy beyond the SMTP status code embedded in it. Common patterns worth branching on:

- `550 5.1.1` — the address doesn't exist. Suppress future sends to it; don't retry.
- `550 5.7.1` — content or policy rejection (often spam filtering). Investigate content, don't just retry blindly.
- `4xx` codes appearing repeatedly across several `message.deferred` events for the same address before an eventual `message.bounced` — usually a full mailbox or a greylisting/rate-limiting pattern on the receiving side, not something wrong with your message.

## Retry behavior you don't need to reimplement

Envelope already retries `message.deferred` outcomes internally with exponential backoff up to 8 attempts before giving up and firing `message.bounced` — don't build your own retry loop that re-sends the same message on a `deferred` event; that just duplicates delivery. Only `message.bounced` is a genuine end state worth reacting to (e.g. suppressing the address, alerting a human, updating your own send status).

## Next steps

- [Receive and verify webhooks](receive-and-verify-webhooks.md)
- [Webhooks → Events](../webhooks/events.md) — full payload reference
