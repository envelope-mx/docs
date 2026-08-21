# API Reference Overview

## Base URL and authentication

The management API is unprefixed — routes are `/accounts`, `/vhosts`, not `/api/v1/...`. Every route except `GET /health` requires a bearer token:

```
Authorization: Bearer env_...
```

See [Authentication](../getting-started/authentication.md) for the token model and [Errors and Responses](../core-concepts/errors-and-responses.md) for the response envelope, status codes, and pagination convention every endpoint below follows.

## Resources

| Resource | Reference |
|---|---|
| Accounts | [accounts.md](accounts.md) |
| Vhosts | [vhosts.md](vhosts.md) |
| Mailboxes | [mailboxes.md](mailboxes.md) |
| Tokens | [tokens.md](tokens.md) |
| Messages (send) | [messages.md](messages.md) |
| Webhooks | [webhooks.md](webhooks.md) |
| Audit log | [audit-log.md](audit-log.md) |
| Data export/delete | [data-export.md](data-export.md) |

## Full route table

| Method | Path | Auth |
|---|---|---|
| <span class="method-badge get">GET</span> | `/health` | none |
| <span class="method-badge post">POST</span> | `/accounts` | admin |
| <span class="method-badge get">GET</span> | `/accounts` | admin |
| <span class="method-badge get">GET</span> | `/accounts/:id` | admin or own account |
| <span class="method-badge get">GET</span> | `/accounts/:accountId/audit-log` | admin or own account |
| <span class="method-badge post">POST</span> | `/accounts/:accountId/vhosts` | admin or own account (self-serve) |
| <span class="method-badge get">GET</span> | `/vhosts` | admin only |
| <span class="method-badge get">GET</span> | `/vhosts/:id` | admin or owning account |
| <span class="method-badge patch">PATCH</span> | `/vhosts/:id/deactivate` | admin or owning account |
| <span class="method-badge patch">PATCH</span> | `/vhosts/:id/policy` | admin or owning account |
| <span class="method-badge post">POST</span> | `/vhosts/:vhostId/mailboxes` | admin or owning account |
| <span class="method-badge get">GET</span> | `/vhosts/:vhostId/mailboxes` | admin or owning account |
| <span class="method-badge get">GET</span> | `/vhosts/:vhostId/mailboxes/:id` | admin or owning account |
| <span class="method-badge delete">DELETE</span> | `/vhosts/:vhostId/mailboxes/:id` | admin or owning account |
| <span class="method-badge post">POST</span> | `/accounts/:accountId/tokens` | admin or own account |
| <span class="method-badge get">GET</span> | `/accounts/:accountId/tokens` | admin or own account |
| <span class="method-badge delete">DELETE</span> | `/accounts/:accountId/tokens/:id` | admin or own account |
| <span class="method-badge post">POST</span> | `/accounts/:accountId/messages` | admin or own account |
| <span class="method-badge post">POST</span> | `/vhosts/:vhostId/webhooks` | admin or owning account |
| <span class="method-badge get">GET</span> | `/vhosts/:vhostId/webhooks` | admin or owning account |
| <span class="method-badge patch">PATCH</span> | `/vhosts/:vhostId/webhooks/:id/disable` | admin or owning account |
| <span class="method-badge get">GET</span> | `/vhosts/:vhostId/webhooks/:id/attempts` | admin or owning account |
| <span class="method-badge post">POST</span> | `/vhosts/:vhostId/webhooks/:id/attempts/:eventId/redrive` | admin or owning account |
| <span class="method-badge get">GET</span> | `/vhosts/:vhostId/audit-log` | admin or owning account |
| <span class="method-badge get">GET</span> | `/vhosts/:id/export` | admin or owning account |
| <span class="method-badge delete">DELETE</span> | `/vhosts/:id/data` | admin or owning account |

<div class="callout gap">
There is no endpoint to list or fetch individual received messages — inbound mail reaches you via a <code>message.received</code> webhook, IMAP, or the bulk per-vhost data export, never a REST "list messages" call. There is also no idempotency-key support on any mutating endpoint — a retried <code>POST</code> after a dropped response is not automatically deduplicated.
</div>

## Next steps

Pick a resource above, or start with [Accounts](accounts.md) if you're following the onboarding order.
