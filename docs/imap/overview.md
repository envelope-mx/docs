# IMAP Overview

Envelope exposes a standard IMAP4rev2 server (`:993` by default, implicit TLS only — not STARTTLS) for reading inbound and quarantined mail with any regular IMAP client, as an alternative or complement to consuming [`message.received` webhooks](../webhooks/events.md).

## Login

`LOGIN <local@vhost> <password>` — the same [mailbox](../api-reference/mailboxes.md) credential [SMTP submission](../smtp/submission.md) uses. There's one shared credential per mailbox for both protocols.

## Mailboxes

Exactly two folders exist, matching where [inbound filtering](../smtp/inbound.md) writes: `INBOX` (accepted mail) and `Quarantine` (flagged-but-not-rejected mail). `SELECT`/`STATUS` work on both by name; `LIST` only ever reports `INBOX`. No folder hierarchy, and no client-created folders — `CREATE`/`RENAME` are not supported.

## Supported

`LOGIN`, `LIST` (INBOX only), `SELECT`, `STATUS`, `FETCH` (`BODY[]`, `FLAGS`, `UID`, `RFC822.SIZE`), `STORE` (flag changes), `NAMESPACE` (a single personal namespace, no prefix), `UNSELECT`.

## Flags

`\Seen`, `\Answered`, `\Flagged`, `\Deleted`, `\Draft` — the complete vocabulary a client may set via `STORE`. `\Recent` and custom keyword flags are not modeled.

## Next steps

- [Limitations](limitations.md) — what's explicitly unsupported and why it matters for client behavior
- [SMTP: Inbound](../smtp/inbound.md) — how mail ends up in these folders in the first place
