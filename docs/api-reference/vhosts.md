# Vhosts

A hosted sending/receiving domain, owned by exactly one account. See [Multi-Tenancy](../core-concepts/multi-tenancy.md) and [DNS and DKIM](../core-concepts/dns-and-dkim.md).

```json
{
  "id": "vh_...", "accountId": "acct_...", "domain": "mail.acme.example", "active": true,
  "maxMessageBytes": 0, "dailyQuota": 0,
  "spamRejectThreshold": 0, "spamQuarantineThreshold": 0,
  "retentionDays": 0,
  "dkimSelector": "envelope",
  "dkimDnsRecord": "v=DKIM1; k=rsa; p=..."
}
```

`dkimDnsRecord` is `omitempty` — only populated by endpoints that hydrate the key (creation and `GET /vhosts/:id`); the list endpoint below never includes it. The four policy fields default to `0`, meaning "unconfigured, platform default applies," not "zero allowed."

## <span class="method-badge post">POST</span> `/accounts/:accountId/vhosts`

Admin, or that account's own token (self-serve). Creates the vhost and generates its DKIM key pair in the same step.

**Request**

```json
{ "domain": "mail.acme.example" }
```

**Response `201`** — the full vhost object above, with `dkimDnsRecord` populated.

Errors: `400` empty domain; `404` account not found; `500` if the domain is already in use by another vhost (a uniqueness conflict that isn't currently mapped to `409` — see [Errors and Responses](../core-concepts/errors-and-responses.md)).

## <span class="method-badge get">GET</span> `/vhosts`

Admin-only. Every vhost across every account on the deployment, cursor-paginated. Does not hydrate `dkimDnsRecord`.

## <span class="method-badge get">GET</span> `/vhosts/:id`

Admin, or the owning account's token. Returns the full vhost object, including `dkimDnsRecord`.

<div class="callout note">
A non-admin token requesting a vhost ID it doesn't own gets <code>403</code> whether that vhost belongs to another account or doesn't exist at all — see <a href="../core-concepts/errors-and-responses.html">Errors and Responses</a> for why.
</div>

## <span class="method-badge patch">PATCH</span> `/vhosts/:id/deactivate`

Admin, or the owning account's token. Stops the vhost accepting new mail; existing data and history remain queryable. No request body. There is no separate "reactivate" endpoint and no way to delete a vhost outright — see [Data Export and Delete](data-export.md) for erasing its stored content.

**Response `200`**: `{ "success": true, "message": "vhost deactivated", "data": null }`

## <span class="method-badge patch">PATCH</span> `/vhosts/:id/policy`

Admin, or the owning account's token. **Full replace, not a merge** — every field below must be supplied, or it overwrites the existing value with its zero value.

**Request**

```json
{
  "maxMessageBytes": 26214400,
  "dailyQuota": 5000,
  "spamRejectThreshold": 15,
  "spamQuarantineThreshold": 6,
  "retentionDays": 90
}
```

**Response `200`** — the updated vhost object.

## Mailboxes

Nested under a vhost — see [Mailboxes](mailboxes.md).

## Next steps

- [DNS and DKIM](../core-concepts/dns-and-dkim.md) — publishing the DKIM record from the create response
- [Mailboxes](mailboxes.md)
- [Set up a new vhost](../guides/set-up-a-new-vhost.md) — the full worked walkthrough
