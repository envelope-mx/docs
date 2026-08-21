# Tokens

Account-scoped bearer credentials — see [Authentication](../getting-started/authentication.md) for the full model.

```json
{ "id": "tok_...", "accountId": "acct_...", "label": "default", "createdAt": "...", "token": "env_..." }
```

`token` is only ever populated on the create response — every other response (list) omits it entirely (`omitempty`), since the raw value is never retrievable again after the moment it's minted. Format: `env_<43-character base64url string>`, persisted only as a SHA-256 hash.

## <span class="method-badge post">POST</span> `/accounts/:accountId/tokens`

Admin, or the account's own token.

**Request**

```json
{ "label": "production-backend" }
```

**Response `201`** — the token object above, with `token` populated. Errors: `404` account not found.

## <span class="method-badge get">GET</span> `/accounts/:accountId/tokens`

Admin, or the account's own token. Cursor-paginated; `token` is omitted from every item.

## <span class="method-badge delete">DELETE</span> `/accounts/:accountId/tokens/:id`

Admin, or the account's own token. Revokes immediately — has no effect on any other token minted for the same account. Errors: `404` if the token doesn't exist.

<div class="callout gap">
Every token minted for an account can act on <strong>every vhost that account owns</strong> — there is no way to mint a token scoped to only one vhost.
</div>

## Next steps

- [Rotate a token](../guides/rotate-a-token.md) — a worked zero-downtime rotation, mint-then-revoke
