# Reverse Proxy: Caddy

The management API has no in-process TLS (see [TLS and Certificates](tls-and-certificates.md)) — Caddy is the lowest-effort way to front it, since it handles automatic HTTPS certificate issuance and renewal on its own, and sets `X-Forwarded-For` by default with no extra configuration.

## Minimal Caddyfile

```
api.yourdomain.example {
    reverse_proxy localhost:8080
}
```

That's the complete configuration for automatic HTTPS plus a correctly-forwarded client IP. Caddy obtains and renews a certificate for `api.yourdomain.example` on its own via Let's Encrypt, and every proxied request already carries `X-Forwarded-For` — nothing else to configure for the [per-IP API rate limit](../core-concepts/rate-limits-and-quotas.md) to work correctly.

## Running it

**As a system service, alongside a [binary](binary.md) deployment:**

```bash
sudo apt install caddy   # or your platform's equivalent
sudo cp Caddyfile /etc/caddy/Caddyfile
sudo systemctl restart caddy
```

**As a container, alongside a [Docker Compose](docker-compose.md) deployment:** see that page's "Adding a reverse proxy to the same file" section — add a `caddy` service pointed at this same `Caddyfile`, mounted read-only.

## Confirm it's working

```bash
curl https://api.yourdomain.example/health
```

If this returns Envelope's health response over a trusted certificate, Caddy is correctly terminating TLS and forwarding to the API. Then confirm the forwarded IP is actually reaching Envelope by triggering the [API rate limit](../core-concepts/rate-limits-and-quotas.md) and checking that a `429` eventually appears for a single misbehaving client, not for everyone.

<div class="callout note">
This Caddyfile only fronts the management API's HTTP port (8080). SMTP (25/587) and IMAP (993) are raw TCP protocols with their own independent TLS story — see <a href="tls-and-certificates.html">TLS and Certificates</a> — Envelope's built-in ACME support handles those directly; there's no need to reverse-proxy them for TLS purposes. Caddy's <code>layer4</code> plugin can front raw TCP if you want a single centralized ingress point for everything, but that's an advanced setup outside what Envelope expects by default.
</div>

## Next steps

- [First Boot and Upgrades](first-boot-and-upgrades.md)
- [Reverse Proxy: nginx](reverse-proxy-nginx.md) — the equivalent if you're standardized on nginx instead
