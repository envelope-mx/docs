# Audit Log

## <span class="method-badge get">GET</span> `/vhosts/:vhostId/audit-log`

Admin, or the owning account's token. Cursor-paginated history of actions taken against this vhost and its child resources (mailboxes, webhooks).

```json
{ "actor": "account:acct_...", "action": "mailbox.create", "detail": "billing", "at": "2026-08-20T12:00:00Z" }
```

| Field | Notes |
|---|---|
| `actor` | `"admin"`, or `"account:<accountId>"` — never a token ID, since a token is just a credential, not the actor of record |
| `action` | See the table below |
| `detail` | Free-text context, `omitempty` |
| `at` | RFC3339 timestamp |

Observed `action` values on this vhost-scoped log: `vhost.create`, `vhost.deactivate`, `vhost.policy.update`, `mailbox.create`, `mailbox.delete`, `webhook.create`, `webhook.disable`, `webhook.redrive`, `message.send`, `vhost.data.export`, `vhost.data.delete`.

<div class="callout gap">
Account-level actions (<code>account.create</code>, <code>token.create</code>, <code>token.revoke</code>) never appear here — their target is an account ID, not a vhost ID, so they're only reachable via <a href="accounts.html"><code>GET /accounts/:accountId/audit-log</code></a> instead. Check both logs if you're auditing everything a given account has done.
</div>

## Next steps

- [Accounts → account-level audit log](accounts.md)
- [Data Export and Delete](data-export.md) — audit history is also included in a full vhost data export
