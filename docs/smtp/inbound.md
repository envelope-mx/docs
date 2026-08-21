# SMTP: Inbound

Envelope receives mail over SMTP on its MX-facing port (`:25` by default). Inbound never requires or accepts authentication — it isn't an `AUTH`-capable session at all.

## Connection

- STARTTLS is offered when TLS is configured, but never required to connect — mirrors how public MX records work in practice.
- A per-connection message size ceiling applies (see [Rate Limits and Quotas](../core-concepts/rate-limits-and-quotas.md) and each vhost's `maxMessageBytes` policy field).

## Per-command behavior

**`MAIL FROM`** — checked against per-source-IP and per-envelope-sender rate limits first (see [Rate Limits and Quotas](../core-concepts/rate-limits-and-quotas.md)). Over either limit → `450 4.7.1` (temporary, retry later).

**`RCPT TO`**:

- Malformed address (no `@`) → `501 5.1.3`
- Domain isn't a registered, active vhost → `550 5.1.1 no such domain here` — rejected **before `DATA` is ever reached**
- Otherwise accepted; the recipient and its resolved vhost are recorded for this transaction

**`DATA`** — read up to the effective size ceiling (the strictest of every recipient vhost's `maxMessageBytes` and the server-wide default); oversized → `552 5.3.4`. The filter pipeline below then runs **once per transaction** (not once per recipient), and the strictest verdict among every recipient vhost governs the single accept/quarantine/reject outcome for the whole transaction — SMTP's `DATA` phase only allows one response, so a message with recipients on two vhosts with different spam thresholds is judged by whichever threshold is stricter.

| Verdict | SMTP response | Storage | Webhook |
|---|---|---|---|
| Reject | `550 5.7.1 message rejected: <reason>` | none | none |
| Quarantine | accepted (200-series) | a distinct `Quarantine` mailbox, not `INBOX` | `message.received` with `quarantine: true` |
| Accept | accepted (200-series) | `INBOX` | `message.received` with `quarantine: false` |

A quarantined message is still accepted at the SMTP level — quarantine is invisible to the sending MTA, visible only to you via the webhook flag or by checking the `Quarantine` folder over [IMAP](../imap/overview.md).

## Filter pipeline

Runs SPF, DKIM signature verification, DMARC alignment, and an rspamd spam score, then decides in this order:

1. Spam score ≥ the vhost's `spamRejectThreshold` (if configured and rspamd was reachable) → **reject**
2. A DMARC record is present, alignment fails, policy is `p=reject` → **reject**
3. rspamd unreachable → **quarantine** (fails open — never blocks mail outright on a scoring-sidecar outage)
4. Spam score ≥ the vhost's `spamQuarantineThreshold` → **quarantine**
5. A DMARC record is present, alignment fails, policy is `p=quarantine` → **quarantine**
6. Otherwise → **accept**

## Next steps

- [SMTP: Submission](submission.md) — the authenticated, outbound-sending counterpart
- [Webhooks → Events](../webhooks/events.md) — the `message.received` payload in full
- [IMAP](../imap/overview.md) — reading what landed in `INBOX` or `Quarantine`
