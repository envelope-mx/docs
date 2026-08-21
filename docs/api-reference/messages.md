# Messages

Send email over REST — a full-featured alternative to [SMTP submission](../smtp/submission.md). Both paths share the same DKIM signing, outbound queue, and `message.queued` webhook.

## <span class="method-badge post">POST</span> `/accounts/:accountId/messages`

Admin, or that account's own token.

**Request**

```json
{
  "from": "Acme Billing <billing@mail.acme.example>",
  "to": ["customer@example.com"],
  "cc": [],
  "bcc": ["archive@acme.example"],
  "subject": "Your invoice",
  "text": "Plain-text body.",
  "html": "<p>HTML body.</p>",
  "attachments": [
    { "filename": "invoice.pdf", "contentType": "application/pdf", "contentBase64": "JVBERi0xLjQK..." }
  ],
  "headers": { "X-Custom-Header": "value" }
}
```

| Field | Type | Notes |
|---|---|---|
| `from` | string | Required. A bare address or `"Display Name <addr>"`. Its domain must resolve to an active vhost owned by this account. |
| `to`, `cc`, `bcc` | string[] | At least one recipient required across all three combined. `bcc` recipients receive the mail but never appear in the composed headers. |
| `subject` | string | |
| `text`, `html` | string | At least one of the two is required — an attachments-only message is rejected. |
| `attachments` | object[] | Each: `filename`, `contentType`, `contentBase64` (standard base64, not base64url). |
| `headers` | object | Custom headers merged into the composed message. May not collide (case-insensitive) with a reserved name — see below. |

**Response `201`**

```json
{ "success": true, "data": { "queued": 1, "vhostId": "vh_...", "jobIds": ["job_..."] } }
```

One job ID per unique resolved recipient (to/cc/bcc deduplicated, order preserved).

## Validation order

Requests are checked in this order — the first failure short-circuits the rest. Authorization is checked **before any resource lookup**, so an unauthorized caller learns nothing about whether the `from` address's vhost even exists.

1. <span class="method-badge delete">403</span> — token doesn't own `:accountId`
2. <span class="method-badge patch">400</span> `from is not a valid email address` — malformed `from`
3. <span class="method-badge delete">404</span> `vhost not found for from address`
4. <span class="method-badge delete">403</span> `vhost is not active`
5. <span class="method-badge delete">403</span> `from address's vhost does not belong to this account` — vhost exists, but under a different account
6. <span class="method-badge delete">503</span> `no DKIM key configured for sending vhost`
7. <span class="method-badge patch">400</span> `at least one recipient (to, cc, or bcc) is required`
8. <span class="method-badge patch">400</span> `text or html body is required`
9. <span class="method-badge patch">400</span> — invalid base64 in an attachment's `contentBase64`
10. <span class="method-badge patch">400</span> `cannot override reserved header "X"` — see below
11. <span class="method-badge patch">429</span> `daily sending quota exceeded for this vhost, try again later` — see [Rate Limits and Quotas](../core-concepts/rate-limits-and-quotas.md)

## Reserved headers

`headers` may not set any of: `from`, `to`, `subject`, `date`, `message-id`, `mime-version`, `content-type`, `content-transfer-encoding` (case-insensitive). These are written by the composer itself — a collision is rejected with `400`, never silently overridden.

<div class="callout gap">
This endpoint enforces the vhost's daily sending quota, but not <code>maxMessageBytes</code> — there is no message-size limit on the REST send path today, unlike SMTP submission. Plan accordingly if size control matters to your integration.
</div>

On success, `message.queued` fires **once per accepted request**, not once per recipient, and a `message.send` entry is recorded to the vhost's [audit log](audit-log.md).

## Next steps

- [Send via REST](../guides/send-via-rest.md) — a full worked example including multi-recipient and attachments
- [Webhooks → Events](../webhooks/events.md) — what happens after the message is queued
- [SMTP: Submission](../smtp/submission.md) — the alternative send path
