# Guide: Receive and Verify Webhooks

End to end: subscribe a webhook, receive a delivery, verify it's genuine, and handle it idempotently.

## 1. Subscribe

```bash
curl -X POST https://mail.yourdomain.example/vhosts/$VHOST_ID/webhooks \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{
    "url": "https://your-app.example/webhooks/envelope",
    "secret": "'"$(openssl rand -hex 32)"'",
    "eventTypes": ["message.received", "message.delivered", "message.bounced", "message.deferred"]
  }'
```

Store the `secret` you generated — the API never returns it again after this response.

## 2. Receive and verify (Node/Express example)

```javascript
const crypto = require("crypto");

const WEBHOOK_SECRET = process.env.ENVELOPE_WEBHOOK_SECRET;
const seenEventIds = new Set(); // use a real store (Redis, DB) in production

app.post("/webhooks/envelope", express.raw({ type: "application/json" }), (req, res) => {
  const signature = req.headers["x-envelope-signature"] || "";
  const expected =
    "sha256=" +
    crypto.createHmac("sha256", WEBHOOK_SECRET).update(req.body).digest("hex");

  const valid =
    signature.length === expected.length &&
    crypto.timingSafeEqual(Buffer.from(signature), Buffer.from(expected));

  if (!valid) {
    return res.status(401).send("invalid signature");
  }

  const event = JSON.parse(req.body);

  if (seenEventIds.has(event.id)) {
    return res.status(200).send("already processed"); // idempotent no-op
  }
  seenEventIds.add(event.id);

  switch (event.type) {
    case "message.received":
      // event.data: { from, to: [...], quarantine }
      break;
    case "message.delivered":
      // event.data: { from, to }
      break;
    case "message.bounced":
      // event.data: { from, to, attempts, error }
      break;
    case "message.deferred":
      // event.data: { from, to, attempts, error, nextAttemptAt }
      break;
  }

  res.status(200).send("ok");
});
```

Two details that matter here: `express.raw()` (not `express.json()`) so the handler sees the **exact bytes** Envelope signed, and `crypto.timingSafeEqual` for a constant-time comparison — see [Signatures](../webhooks/signatures.md) for why both matter.

## 3. Handle retries gracefully

If your endpoint is briefly down or slow, Envelope retries with backoff for up to 8 attempts over roughly 30 minutes before dead-lettering the delivery (see [Webhooks Overview](../webhooks/overview.md)). Returning anything other than `2xx` — including a `500` from a bug in your handler — triggers this same retry behavior, so a transient error on your end is usually self-healing without any action from you. If a delivery does dead-letter, you can inspect and manually [redrive](../api-reference/webhooks.md) it once your endpoint is healthy again.

## Next steps

- [Handle bounces and deferrals](handle-bounces-and-deferrals.md)
- [Webhooks → Events](../webhooks/events.md) — full payload reference for every event type
