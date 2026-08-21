# Webhooks

Endpoints for managing webhook subscriptions and inspecting delivery attempts. For the event catalog and payload shapes themselves, see [Webhooks → Events](../webhooks/events.md); for verifying delivery authenticity, see [Webhooks → Signatures](../webhooks/signatures.md).

```json
{ "id": "wh_...", "vhostId": "vh_...", "eventTypes": ["message.received", "message.bounced"], "disabled": false }
```

`url` and `secret` are write-only — never echoed back in any response, including the create response.

## <span class="method-badge post">POST</span> `/vhosts/:vhostId/webhooks`

Admin, or the owning account's token.

**Request**

```json
{
  "url": "https://your-app.example/webhooks/envelope",
  "secret": "a shared secret you generate",
  "eventTypes": ["message.received", "message.delivered", "message.bounced"]
}
```

An empty or omitted `eventTypes` subscribes to every event type. **Response `201`** — the subscription object above. Errors: `400` empty `url`/`secret`; `404` vhost not found.

## <span class="method-badge get">GET</span> `/vhosts/:vhostId/webhooks`

Admin, or the owning account's token. Cursor-paginated.

## <span class="method-badge patch">PATCH</span> `/vhosts/:vhostId/webhooks/:id/disable`

Admin, or the owning account's token. Soft-disable — delivery history is preserved and still queryable via the attempts endpoint below. No request body.

<div class="callout gap">
There's no re-enable endpoint and no way to update a subscription's URL or event types after creation — disable it and create a new one instead.
</div>

## <span class="method-badge get">GET</span> `/vhosts/:vhostId/webhooks/:id/attempts`

Admin, or the owning account's token. Delivery attempt history for one subscription, oldest first, cursor-paginated.

```json
{ "id": "att_...", "eventId": "evt_...", "attempt": 3, "statusCode": 0, "error": "connection refused", "attemptedAt": "..." }
```

## <span class="method-badge post">POST</span> `/vhosts/:vhostId/webhooks/:id/attempts/:eventId/redrive`

Admin, or the owning account's token. Manually retries a specific dead-lettered event against this subscription, using its **current** URL and secret (not whatever they were at the time of the original attempt). No request body.

**Response `200`**: `{ "success": true, "message": "webhook delivery redriven" }`. Errors: `404` if the subscription doesn't exist under this vhost, or `eventId` was never attempted against it; `503` if this deployment never wired a dispatcher into the API (an operator-side configuration gap, not something a caller can work around).

## Next steps

- [Webhooks → Events](../webhooks/events.md) — the full event and payload catalog
- [Webhooks → Signatures](../webhooks/signatures.md) — verifying `X-Envelope-Signature`
- [Receive and verify webhooks](../guides/receive-and-verify-webhooks.md) — a worked end-to-end guide
