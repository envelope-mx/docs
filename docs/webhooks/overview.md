# Webhooks Overview

Subscribe a URL to a vhost via [`POST /vhosts/:vhostId/webhooks`](../api-reference/webhooks.md) and Envelope will `POST` a signed JSON envelope to it whenever a matching event fires.

## Delivery mechanics

- `POST`, `Content-Type: application/json`, body is the [event envelope](events.md).
- `X-Envelope-Signature: sha256=<hex>` — an HMAC-SHA256 of the raw request body, keyed with the subscription's secret. See [Signatures](signatures.md) for verifying it.
- `X-Envelope-Event: <event type>` — the same value as the body's `type` field, for routing without parsing JSON first if you prefer.
- Any `2xx` response is treated as success. The response body is read (up to 4096 bytes) but never inspected — return whatever you like, just return `2xx`.

## Retry and dead-letter

A non-`2xx` response (or no response at all) is retried with full-jitter exponential backoff: 30 second base, 30 minute cap, 8 attempts by default. Exhausting every attempt **dead-letters** the delivery — it stops retrying automatically, but the complete attempt history stays queryable forever via [`GET .../webhooks/:id/attempts`](../api-reference/webhooks.md), and you can manually trigger one more try with the [redrive endpoint](../api-reference/webhooks.md) at any time afterward.

## Fan-out

One event fans out to every enabled subscription on its vhost whose `eventTypes` is either empty (matches everything) or explicitly lists that event type. Disabled subscriptions are skipped entirely — no attempt is even recorded against them.

## Why this can never block mail flow

Firing a webhook is a fast, durable local database insert — the actual outbound HTTP `POST` happens later, from a background dispatch loop. A slow or completely unreachable subscriber endpoint can never block SMTP acceptance or an API response; at worst, its own deliveries pile up in the retry queue and eventually dead-letter.

## Next steps

- [Events](events.md) — the full event and payload catalog
- [Signatures](signatures.md) — verifying `X-Envelope-Signature`
- [Receive and verify webhooks](../guides/receive-and-verify-webhooks.md) — a worked end-to-end guide
