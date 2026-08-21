# Errors and Responses

## The response envelope

Every management API response, success or error, has the same JSON shape:

```json
{
  "success": true,
  "data": { },
  "message": "optional — omitted entirely when empty",
  "meta": { "optional": "omitted entirely when nil, used for pagination" }
}
```

`Content-Type: application/json; charset=utf-8` on every response.

## Status codes

| Code | Meaning | `success` | Seen on |
|---|---|---|---|
| <span class="method-badge get">200</span> | OK | `true` | successful `GET`, and mutations that don't create a resource (`PATCH`/`DELETE`) |
| <span class="method-badge post">201</span> | Created | `true` | every `POST` that creates a resource |
| <span class="method-badge patch">400</span> | Bad Request | `false` | missing/malformed field, reserved-header collision, invalid attachment encoding |
| <span class="method-badge patch">401</span> | Unauthorized | `false` | missing/malformed `Authorization` header, unrecognized or revoked token |
| <span class="method-badge delete">403</span> | Forbidden | `false` | admin required but caller isn't; token doesn't own the target account/vhost |
| <span class="method-badge delete">404</span> | Not Found | `false` | resource genuinely missing (admin caller, or an account-level lookup) |
| <span class="method-badge patch">429</span> | Too Many Requests | `false` | API per-IP rate limit; REST send daily quota exceeded |
| <span class="method-badge delete">500</span> | Internal Server Error | `false` | any unhandled failure — **including duplicate-domain/duplicate-mailbox conflicts**, see below |
| <span class="method-badge delete">503</span> | Service Unavailable | `false` | webhook redrive with no dispatcher configured; REST send to a vhost with no DKIM key |

<div class="callout gap">
Uniqueness conflicts (creating a vhost for a domain that already exists, or a mailbox local-part that already exists on a vhost) currently surface as a generic <code>500</code>, not a REST-idiomatic <code>409 Conflict</code>. Treat a <code>500</code> immediately after a create call as a possible duplicate, not only as a server fault.
</div>

## Why a foreign or missing vhost ID both return 403

A non-admin token's request against a `:vhostId` it doesn't own returns the identical `403 "token is not authorized for this vhost"` whether that vhost belongs to a different account **or doesn't exist at all**. This is deliberate: distinguishing the two with a `404` would let any authenticated account enumerate which vhost IDs exist across every other tenant on the deployment. Admin callers are unaffected — they get a genuine `404` when a vhost is actually missing.

## Pagination

Every list endpoint (`GET /accounts`, `GET /vhosts`, mailboxes, tokens, webhook subscriptions, webhook attempts, audit logs) is cursor-paginated via `?cursor=&limit=` query parameters:

```json
{ "success": true, "data": [ /* ... */ ], "meta": { "nextCursor": "<id, or empty when this is the last page>" } }
```

Pass the previous response's `meta.nextCursor` as the next request's `cursor` to continue. An empty `nextCursor` means you've reached the end. Default and maximum `limit` vary slightly by resource (100/500 for most; see each resource's page in the [API Reference](../api-reference/overview.md)).

## Next steps

- [Rate Limits and Quotas](rate-limits-and-quotas.md)
- [API Reference Overview](../api-reference/overview.md)
