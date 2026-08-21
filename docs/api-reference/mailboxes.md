# Mailboxes

A send/receive identity on a vhost — what SMTP `AUTH` and IMAP `LOGIN` authenticate against. Only needed for [SMTP submission](../smtp/submission.md) or [IMAP](../imap/overview.md); REST send and webhooks don't require one.

```json
{ "id": "mb_...", "vhostId": "vh_...", "localPart": "billing" }
```

The password is never returned — only `id`, `vhostId`, and `localPart`.

## <span class="method-badge post">POST</span> `/vhosts/:vhostId/mailboxes`

Admin, or the owning account's token.

**Request**

```json
{ "localPart": "billing", "password": "a strong password" }
```

**Response `201`** — the mailbox object above. Errors: `400` empty `localPart`/`password`; `404` vhost not found; `500` if that `localPart` already exists on this vhost (uniqueness conflict, see [Errors and Responses](../core-concepts/errors-and-responses.md)).

## <span class="method-badge get">GET</span> `/vhosts/:vhostId/mailboxes`

Admin, or the owning account's token. Cursor-paginated list.

## <span class="method-badge get">GET</span> `/vhosts/:vhostId/mailboxes/:id`

Admin, or the owning account's token. Single mailbox.

## <span class="method-badge delete">DELETE</span> `/vhosts/:vhostId/mailboxes/:id`

Admin, or the owning account's token. No response body beyond the success envelope.

<div class="callout gap">
There is no endpoint to change a mailbox's password after creation — deleting and recreating it is the only way to rotate a mailbox credential today.
</div>

## Next steps

- [SMTP: Submission](../smtp/submission.md) — using a mailbox's credential to send
- [IMAP](../imap/overview.md) — using it to read mail
