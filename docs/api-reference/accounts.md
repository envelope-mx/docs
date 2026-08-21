# Accounts

An account is a business or tenant. See [Multi-Tenancy](../core-concepts/multi-tenancy.md) for the full ownership model.

```json
{ "id": "acct_...", "name": "Acme Inc", "createdAt": "2026-08-20T12:00:00Z" }
```

## <span class="method-badge post">POST</span> `/accounts`

Admin-only. Creates a new account and immediately mints its first token in the same response.

**Request**

```json
{ "name": "Acme Inc" }
```

**Response `201`**

```json
{
  "success": true,
  "data": {
    "account": { "id": "acct_...", "name": "Acme Inc", "createdAt": "2026-08-20T12:00:00Z" },
    "token": { "id": "tok_...", "accountId": "acct_...", "label": "default", "createdAt": "2026-08-20T12:00:00Z", "token": "env_..." }
  }
}
```

`data.token.token` is the raw bearer token, visible **only in this response** — store it now. Errors: `400` if `name` is empty.

## <span class="method-badge get">GET</span> `/accounts`

Admin-only. Cursor-paginated list of every account on the deployment.

```
GET /accounts?cursor=&limit=100
```

**Response `200`**

```json
{ "success": true, "data": [ { "id": "acct_...", "name": "Acme Inc", "createdAt": "..." } ], "meta": { "nextCursor": "" } }
```

## <span class="method-badge get">GET</span> `/accounts/:id`

Admin, or that account's own token. Returns a single account.

## <span class="method-badge get">GET</span> `/accounts/:accountId/audit-log`

Admin, or that account's own token. Cursor-paginated audit entries scoped to this account — `account.create`, `token.create`, `token.revoke`, and any other account-level action, none of which are reachable through the vhost-scoped [audit log](audit-log.md) since their target is an account ID, not a vhost ID.

```json
{ "success": true, "data": [ { "actor": "admin", "action": "account.create", "detail": "", "at": "..." } ], "meta": { "nextCursor": "" } }
```

<div class="callout gap">
There is no self-serve account creation and no way to delete or rename an account via the API — <code>POST /accounts</code> is always admin-gated, and there's no <code>PATCH</code>/<code>DELETE</code> for the resource at all.
</div>

## Next steps

- [Vhosts](vhosts.md) — the next step after creating an account
- [Tokens](tokens.md) — minting further tokens for this account
