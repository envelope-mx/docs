# Guide: Set Up a New Vhost

A production-ready vhost setup, from account to first real send — ties together [Multi-Tenancy](../core-concepts/multi-tenancy.md), [DNS and DKIM](../core-concepts/dns-and-dkim.md), and [Rate Limits and Quotas](../core-concepts/rate-limits-and-quotas.md) into one walkthrough. If you just want the fastest possible path, see the [Quickstart](../getting-started/quickstart.md) instead — this guide additionally covers DNS propagation and real policy limits before you'd actually rely on this in production.

## 1. Create the vhost

```bash
curl -X POST https://mail.yourdomain.example/accounts/$ACCOUNT_ID/vhosts \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"domain": "mail.acme.example"}'
```

Save `data.dkimDnsRecord` and `data.id` (as `$VHOST_ID`) from the response.

## 2. Publish DNS records

```
envelope._domainkey.mail.acme.example.   TXT   "v=DKIM1; k=rsa; p=..."   (from step 1's dkimDnsRecord)
mail.acme.example.                       TXT   "v=spf1 ip4:<your deliverer's egress IP> ~all"
_dmarc.mail.acme.example.                TXT   "v=DMARC1; p=quarantine; rua=mailto:dmarc-reports@acme.example"
```

If you're also receiving inbound mail on this domain, ask your operator for the `smtp-inbound` host and add an MX record pointing at it. See [DNS and DKIM](../core-concepts/dns-and-dkim.md) for what Envelope does and doesn't generate for you.

Wait for propagation before relying on real sending — check with `dig TXT envelope._domainkey.mail.acme.example` (or your provider's DNS tools) until the record resolves. There's no readiness endpoint that confirms this for you.

## 3. Set real policy limits

A freshly-created vhost has every policy field at `0` (unconfigured — a platform default applies). Set real values before production traffic:

```bash
curl -X PATCH https://mail.yourdomain.example/vhosts/$VHOST_ID/policy \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{
    "maxMessageBytes": 26214400,
    "dailyQuota": 5000,
    "spamRejectThreshold": 15,
    "spamQuarantineThreshold": 6,
    "retentionDays": 90
  }'
```

Remember this is a **full replace** — include every field, not just the ones you're changing.

## 4. (Optional) Create a mailbox for SMTP/IMAP

Skip this if you're only using [REST send](../api-reference/messages.md) and webhooks.

```bash
curl -X POST https://mail.yourdomain.example/vhosts/$VHOST_ID/mailboxes \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"localPart": "billing", "password": "a strong password"}'
```

## 5. Subscribe to webhooks

```bash
curl -X POST https://mail.yourdomain.example/vhosts/$VHOST_ID/webhooks \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"url": "https://your-app.example/webhooks/envelope", "secret": "'"$(openssl rand -hex 32)"'"}'
```

## 6. Send a real test

```bash
curl -X POST https://mail.yourdomain.example/accounts/$ACCOUNT_ID/messages \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"from": "billing@mail.acme.example", "to": ["you@example.com"], "subject": "Live test", "text": "Real DNS, real policy, real send."}'
```

Confirm `message.queued` fires, then `message.delivered` shortly after — if you instead see `message.bounced` or `message.deferred`, check `data.error` against [Handle Bounces and Deferrals](handle-bounces-and-deferrals.md), and double-check the DKIM/SPF records from step 2 have actually propagated.

## Next steps

- [Send via REST](send-via-rest.md) or [Send via SMTP](send-via-smtp.md) for your real integration
- [Receive and verify webhooks](receive-and-verify-webhooks.md)
- [Rotate a token](rotate-a-token.md) — once this is in production, plan for credential rotation
