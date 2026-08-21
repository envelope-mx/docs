# Quickstart

This walks through sending your first email against an already-running Envelope instance, using only `curl`. If you don't have an instance running yet, see [Deployment](../deployment/overview.md) first.

Every request below assumes your instance's management API is reachable at `https://mail.yourdomain.example` — substitute your own host.

## 1. Create an account

This is the **only** step in the entire flow that needs the admin bootstrap token (`ENVELOPE_API_ADMIN_TOKEN`, see [Deployment Overview](../deployment/overview.md)). Everything after this uses the account's own token.

```bash
curl -X POST https://mail.yourdomain.example/accounts \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Acme Inc"}'
```

```json
{
  "success": true,
  "data": {
    "account": { "id": "acct_...", "name": "Acme Inc", "createdAt": "2026-08-20T12:00:00Z" },
    "token": { "id": "tok_...", "accountId": "acct_...", "label": "default", "createdAt": "2026-08-20T12:00:00Z", "token": "env_..." }
  }
}
```

Save `data.token.token` as `$ACCOUNT_TOKEN` and `data.account.id` as `$ACCOUNT_ID`. The raw token value is shown **exactly once** — it's stored hashed and can never be retrieved again. If you lose it, revoke it and [mint a new one](../api-reference/tokens.md).

## 2. Create a vhost (self-serve, no admin token needed)

```bash
curl -X POST https://mail.yourdomain.example/accounts/$ACCOUNT_ID/vhosts \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"domain": "mail.acme.example"}'
```

```json
{
  "success": true,
  "data": {
    "id": "vh_...", "accountId": "acct_...", "domain": "mail.acme.example", "active": true,
    "maxMessageBytes": 0, "dailyQuota": 0, "spamRejectThreshold": 0, "spamQuarantineThreshold": 0,
    "retentionDays": 0, "dkimSelector": "envelope",
    "dkimDnsRecord": "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC..."
  }
}
```

A 2048-bit DKIM key was generated for this vhost automatically. Publish `data.dkimDnsRecord` as a TXT record at `envelope._domainkey.mail.acme.example` — see [DNS and DKIM](../core-concepts/dns-and-dkim.md) for the full record and what else (MX, SPF, DMARC) you're responsible for.

The four policy fields (`maxMessageBytes`, `dailyQuota`, the two spam thresholds, `retentionDays`) all default to `0`, meaning "unconfigured — platform default applies." Set them with [`PATCH /vhosts/:id/policy`](../api-reference/vhosts.md) once you know your real limits.

## 3. Send your first email

```bash
curl -X POST https://mail.yourdomain.example/accounts/$ACCOUNT_ID/messages \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "from": "hello@mail.acme.example",
    "to": ["you@example.com"],
    "subject": "Hello from Envelope",
    "text": "It works.",
    "html": "<p>It works.</p>"
  }'
```

```json
{ "success": true, "data": { "queued": 1, "vhostId": "vh_...", "jobIds": ["job_..."] } }
```

The message is DKIM-signed, staged, and queued for delivery — the same pipeline authenticated SMTP submission uses. A `message.queued` webhook fires once for the request (not once per recipient); if you [subscribe a webhook](../api-reference/webhooks.md) you'll also see `message.delivered` (or `message.bounced`/`message.deferred`) once the deliverer actually reaches the destination MX. Full field reference: [API Reference → Messages](../api-reference/messages.md).

<div class="callout warning">
Until the DKIM DNS record from step 2 has propagated, outbound mail from this vhost may be signed with a key receiving mail servers can't yet verify — some providers will still accept it, others may reject or spam-foldering it. Publish the record before sending anything that needs to reliably land.
</div>

## Next steps

- [Authentication](authentication.md) — the full token model and self-serve onboarding sequence
- [Multi-Tenancy](../core-concepts/multi-tenancy.md) — how accounts, vhosts, and mailboxes relate
- [Send via SMTP](../guides/send-via-smtp.md) — the alternative to REST send, using existing mail-sending code
- [Receive and verify webhooks](../guides/receive-and-verify-webhooks.md) — reacting to inbound mail and delivery events
