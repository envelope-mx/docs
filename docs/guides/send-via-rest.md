# Guide: Send via REST

A complete example sending to multiple recipients with an HTML+text body and an attachment, using the [REST send endpoint](../api-reference/messages.md).

## Prerequisites

An account token and an active vhost with a DKIM key — see the [Quickstart](../getting-started/quickstart.md) if you don't have these yet.

## Build the request

```bash
ATTACHMENT_B64=$(base64 -i invoice.pdf | tr -d '\n')

curl -X POST https://mail.yourdomain.example/accounts/$ACCOUNT_ID/messages \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" \
  -H "Content-Type: application/json" \
  -d @- <<EOF
{
  "from": "Acme Billing <billing@mail.acme.example>",
  "to": ["customer@example.com"],
  "cc": ["accounting@acme.example"],
  "subject": "Your invoice for August",
  "text": "Hi, your invoice is attached. Thanks for your business.",
  "html": "<p>Hi,</p><p>Your invoice is attached. Thanks for your business.</p>",
  "attachments": [
    { "filename": "invoice.pdf", "contentType": "application/pdf", "contentBase64": "$ATTACHMENT_B64" }
  ],
  "headers": { "X-Invoice-Id": "INV-1042" }
}
EOF
```

```json
{ "success": true, "data": { "queued": 2, "vhostId": "vh_...", "jobIds": ["job_a...", "job_b..."] } }
```

`queued: 2` because `to` and `cc` together resolved to two unique recipients — each gets its own queue job and its own eventual `message.delivered`/`bounced`/`deferred` webhook, but only **one** `message.queued` event fires for the whole request.

## What to check next

1. Subscribe a webhook (see [API Reference → Webhooks](../api-reference/webhooks.md)) if you haven't, and confirm `message.queued` arrives with `data.to` containing both addresses.
2. Watch for `message.delivered` per recipient. A `message.bounced` or `message.deferred` instead means something went wrong at the destination — see [Handle bounces and deferrals](handle-bounces-and-deferrals.md).

## Common mistakes

- **`from`'s domain isn't a vhost you own** → `403 from address's vhost does not belong to this account`. The `from` address must resolve to a vhost owned by the same account as `$ACCOUNT_TOKEN`.
- **Neither `text` nor `html` set** → `400 text or html body is required`. Attachments alone aren't a valid message body.
- **A custom header collides with a reserved name** (`from`, `to`, `subject`, `date`, `message-id`, `mime-version`, `content-type`, `content-transfer-encoding`) → `400`. Pick a different header name.
- **Large attachments** — there's no message-size limit enforced on this endpoint specifically (unlike SMTP submission); very large attachments will still cost you on the receiving end's own size limits, so keep them reasonable regardless.

## Next steps

- [Send via SMTP](send-via-smtp.md) — the alternative path
- [Receive and verify webhooks](receive-and-verify-webhooks.md)
