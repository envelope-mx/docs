# DNS and DKIM

## What Envelope generates for you

The moment a vhost is created (`POST /accounts/:accountId/vhosts`), Envelope generates a **2048-bit RSA DKIM key pair** for it server-side — you never supply or see the private key. The selector is always the literal string `envelope` and is **not configurable** per vhost.

The public key is returned as a ready-to-publish DNS TXT record value in the vhost response's `dkimDnsRecord` field:

```
v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC...
```

Publish it as a TXT record at:

```
envelope._domainkey.mail.acme.example
```

`dkimDnsRecord` is only populated by endpoints that hydrate the key: `POST .../vhosts` (creation) and `GET /vhosts/:id` both include it; the cross-tenant `GET /vhosts` list does not. See [API Reference → Vhosts](../api-reference/vhosts.md).

<div class="callout note">
A 2048-bit RSA public key's base64 encoding is long enough that some DNS providers require splitting it across multiple quoted strings within the same TXT record (a standard DNS TXT record mechanic, not an Envelope-specific concern) — check your provider's TXT record length limit if publication fails.
</div>

There is no key-rotation endpoint in the API — DKIM key rotation is an operator-run maintenance operation, not something a tenant triggers via a request.

## What you're responsible for

Envelope does not generate, publish, or verify SPF or DMARC records — there's no field or endpoint for either. You'll typically want:

```
# SPF — authorize Envelope's outbound egress IP(s) to send as this domain
mail.acme.example.   TXT   "v=spf1 ip4:<your deliverer's egress IP> ~all"

# DMARC — start at "quarantine" while validating, tighten to "reject" later
_dmarc.mail.acme.example.   TXT   "v=DMARC1; p=quarantine; rua=mailto:dmarc-reports@acme.example"
```

If you're receiving inbound mail on this vhost, your operator also needs an **MX record** pointing this domain at their `smtp-inbound` host — that hostname isn't returned by any API call, since it's deployment-specific; ask whoever operates your Envelope instance.

<div class="callout gap">
There's no readiness or DNS-propagation-check endpoint. <code>GET /vhosts/:id</code> only reflects Envelope's own internal state, never whether your DNS records have actually propagated or resolve correctly — verify with an external DNS lookup tool before relying on real deliverability.
</div>

## Next steps

- [Rate Limits and Quotas](rate-limits-and-quotas.md)
- [Set up a new vhost](../guides/set-up-a-new-vhost.md) — the full worked walkthrough from account to first send
- [API Reference → Vhosts](../api-reference/vhosts.md)
