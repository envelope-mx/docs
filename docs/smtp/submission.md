# SMTP: Submission

Authenticated SMTP send — the traditional alternative to [REST send](../api-reference/messages.md), for integrating existing mail-sending code (`smtplib`, `nodemailer`, an MTA relay configuration, etc.) without changing it. Port `587` by default.

## Authentication

- Only `AUTH PLAIN` is supported — no `LOGIN`, no OAuth2 bearer.
- **The username must be a full email address** (`local@vhost`) — this is how the session learns which vhost's DKIM key to sign with. It's authenticated against a [mailbox's](../api-reference/mailboxes.md) bcrypt-hashed password, the same credential [IMAP](../imap/overview.md) uses.
- `AUTH` requires a prior `STARTTLS` in every real deployment — plaintext `AUTH` is a dev/test-only allowance, never enabled in production.
- Any command before authentication succeeds, including `MAIL FROM`, is rejected: `530 5.7.0 authentication required`.

## Sending

Once authenticated, `MAIL FROM`/`RCPT TO`/`DATA` proceed as normal SMTP. On `DATA`:

1. No DKIM key configured for the authenticated vhost → `451 no DKIM key configured for sending vhost`
2. Daily quota check (see [Rate Limits and Quotas](../core-concepts/rate-limits-and-quotas.md)) → `452 4.7.1 daily sending quota exceeded for this vhost, try again later`
3. The raw body is DKIM-signed automatically with the authenticated vhost's key and selector — the sender never handles key material directly
4. Signing failure → `451 DKIM signing failed`
5. The signed message is staged and one durable queue job is enqueued per recipient — storage/queue failure → `451 temporary storage failure` / `451 queue unavailable`

On success, `message.queued` fires **once per accepted submission**, not once per recipient — identical semantics to the REST send endpoint.

<div class="callout note">
A mailbox's credential can only ever send DKIM-signed mail as its own vhost's domain — there's no way to authenticate as one vhost and send <code>From:</code> a different one, regardless of what the message body's <code>From:</code> header claims.
</div>

## Next steps

- [API Reference → Messages](../api-reference/messages.md) — the REST alternative, same underlying pipeline
- [Send via SMTP](../guides/send-via-smtp.md) — a worked example
- [Webhooks → Events](../webhooks/events.md)
