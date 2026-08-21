# Data Export and Delete

Compliance-oriented endpoints for exporting or erasing a vhost's stored data.

## <span class="method-badge get">GET</span> `/vhosts/:id/export`

Admin, or the owning account's token. A complete snapshot: the vhost record, its mailboxes, every stored message (base64-encoded body included), its webhook subscriptions, and its audit log — **not paginated**, the entire export is buffered into one response.

```json
{
  "vhost": { "...": "vhostView" },
  "mailboxes": [ { "...": "mailboxView" } ],
  "messages": [
    { "mailbox": "billing@mail.acme.example", "size": 4213, "flags": ["\\Seen"], "createdAt": "...", "body": "<base64>" }
  ],
  "webhookSubscriptions": [ { "...": "webhookSubscriptionView" } ],
  "auditLog": [ { "...": "auditEntryView" } ]
}
```

Best-effort on messages: one that's deleted by a concurrent retention purge between being listed and being read is silently skipped, not treated as an error.

## <span class="method-badge delete">DELETE</span> `/vhosts/:id/data`

Admin, or the owning account's token. Erases stored message **content only** — mailbox credentials, the vhost record, and the account all persist. Idempotent: deleting an already-empty vhost returns `messagesDeleted: 0`, not an error.

**Response `200`**

```json
{ "success": true, "message": "vhost data deleted", "data": { "messagesDeleted": 42 } }
```

<div class="callout gap">
There is no endpoint to delete a vhost or account outright — only <a href="vhosts.html"><code>PATCH /vhosts/:id/deactivate</code></a> (stop it accepting mail) and this content-erase endpoint. The vhost and account records themselves are permanent.
</div>

## Next steps

- [Audit Log](audit-log.md)
- [Rate Limits and Quotas](../core-concepts/rate-limits-and-quotas.md) — retention-window interaction with `retentionDays`
