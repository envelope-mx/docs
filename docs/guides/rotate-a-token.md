# Guide: Rotate a Token

Zero-downtime credential rotation — mint the replacement, switch your integration over, then revoke the old one. Never revoke first; you'll lock yourself out of the account until you generate a new token, which itself needs the token you just revoked (or the admin token) to do.

## 1. Mint a new token

```bash
curl -X POST https://mail.yourdomain.example/accounts/$ACCOUNT_ID/tokens \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"label": "production-backend-2026-08"}'
```

```json
{ "success": true, "data": { "id": "tok_new...", "accountId": "acct_...", "label": "production-backend-2026-08", "createdAt": "...", "token": "env_..." } }
```

Copy `data.token` immediately — it's shown exactly once. Give the new token a `label` that identifies what will use it and roughly when it was minted; you'll thank yourself when deciding what's safe to revoke later.

## 2. Roll out the new token

Update whatever holds the old credential (a secrets manager, an environment variable, a config file) with the new token, and deploy. Both the old and new tokens work simultaneously during this window — minting a new token has no effect on any existing one.

## 3. Confirm the new token is actually in use

Check your application's logs or metrics for a request made with the new token before proceeding — a `401` on the new credential at this point means something in the rollout didn't take, and revoking the old token next would cause an outage.

## 4. List tokens and revoke the old one

```bash
curl https://mail.yourdomain.example/accounts/$ACCOUNT_ID/tokens \
  -H "Authorization: Bearer $ACCOUNT_TOKEN"
```

```json
{ "success": true, "data": [
  { "id": "tok_old...", "accountId": "acct_...", "label": "production-backend-2025-11", "createdAt": "..." },
  { "id": "tok_new...", "accountId": "acct_...", "label": "production-backend-2026-08", "createdAt": "..." }
], "meta": { "nextCursor": "" } }
```

```bash
curl -X DELETE https://mail.yourdomain.example/accounts/$ACCOUNT_ID/tokens/tok_old... \
  -H "Authorization: Bearer $ACCOUNT_TOKEN"
```

Revoking `tok_old...` has no effect on `tok_new...` or any other token minted for this account.

<div class="callout note">
Rotation is per-token, not per-account — every token you mint for an account can act on every vhost that account owns (see <a href="../getting-started/authentication.html">Authentication</a>), so rotating one credential used by one service doesn't require touching any other service's token.
</div>

## Next steps

- [API Reference → Tokens](../api-reference/tokens.md)
- [Authentication](../getting-started/authentication.md)
