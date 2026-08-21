# Rate Limits and Quotas

Every limit below shares the same shape — a token-bucket with a capacity and a per-second refill rate — and every one is independently nil-able on the operator's side: if a deployment hasn't wired a given dimension up, that dimension simply doesn't apply, rather than erroring. Token authentication itself is unaffected either way.

| Dimension | Applies to | Default capacity / refill | Over-limit response |
|---|---|---|---|
| Per-source-IP | Inbound SMTP (`MAIL FROM`) | 20 / 1 per sec | SMTP `450 4.7.1` (temporary — retry) |
| Per-envelope-sender | Inbound SMTP (`MAIL FROM`) | 10 / 0.5 per sec | SMTP `450 4.7.1` (temporary — retry) |
| Per-vhost daily sending quota | SMTP submission **and** REST send — identical enforcement on both paths | `dailyQuota` policy field (`0` = disabled) | SMTP `452 4.7.1`; REST <span class="method-badge post">429</span> |
| Management API per-client-IP | Every `internal` management API request, checked **before** the token itself is even looked up | 60 / 1 per sec | <span class="method-badge post">429</span> |

## The daily sending quota specifically

Set per-vhost via [`PATCH /vhosts/:id/policy`](../api-reference/vhosts.md)'s `dailyQuota` field — an integer message count. It's enforced as a rolling-window token bucket (capacity = `dailyQuota`, refill = `dailyQuota / 86400` per second), not a hard reset at midnight. A vhost with `dailyQuota: 0` (the default for a freshly-created vhost) has **no** quota enforced at all.

- SMTP submission over quota → `452 4.7.1 daily sending quota exceeded for this vhost, try again later`
- REST send over quota → `429` with the same message text in the response body

<div class="callout gap">
The REST send endpoint enforces the daily quota but does <strong>not</strong> enforce <code>maxMessageBytes</code> — there's no message-size limit on that path today, unlike SMTP submission which rejects an oversized message mid-transfer. Keep this in mind if you're relying on the size policy for cost or abuse control and using REST send.
</div>

## The management API's per-IP limit needs a reverse proxy

The client IP for this limiter is read from the `X-Forwarded-For` header — never the raw TCP connection — since the management API has no direct visibility into the connecting socket. If nothing in front of Envelope sets that header, this dimension silently does nothing (every caller is treated as unrateable-by-IP; token auth still protects every endpoint regardless). Both [Caddy](../deployment/reverse-proxy-caddy.md) and [nginx](../deployment/reverse-proxy-nginx.md) reverse-proxy guides show the header configuration this needs.

No `Retry-After` header is set on any `429`/`452`/`450` response — just the status code and a plain-text message. Build your own backoff.

## Next steps

- [Errors and Responses](errors-and-responses.md)
- [API Reference → Vhosts](../api-reference/vhosts.md) — setting `dailyQuota` and the other policy fields
- [Deployment Overview](../deployment/overview.md) — the operator-side environment variables behind the IP-based limits
