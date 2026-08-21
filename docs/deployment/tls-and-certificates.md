# TLS and Certificates

Envelope's TLS story has two completely separate halves: the mail protocols (SMTP, IMAP) handle their own TLS in-process; the management API does not handle TLS at all.

## SMTP and IMAP

- **Default**: a self-signed certificate, generated per-process. Fine for local testing and internal use; browsers and mail clients won't trust it, and neither will remote MTAs verifying inbound TLS on delivery.
- **Real certificates**: set `ENVELOPE_ACME_EMAIL` (and a real `ENVELOPE_DOMAIN`) on every `smtp-inbound`/`smtp-submission`/`imap` process. Certificates are then issued automatically on first handshake, **per active vhost domain specifically** — issuance is scoped to domains actually registered as active vhosts, which is what prevents anyone from burning your ACME rate-limit budget requesting certificates for arbitrary hostnames. Needs port 80 reachable from the public internet (the ACME HTTP-01 challenge).
- Setting `ENVELOPE_ACME_EMAIL` is your explicit agreement to the CA's subscriber agreement — there's no separate prompt.
- Validate with `ENVELOPE_ACME_STAGING=1` first: exercises the real issuance flow, produces untrusted certificates, and doesn't spend your production rate-limit budget. Remove it once you've confirmed issuance works.
- Certificate/account state persists to local storage inside the container/host (`~/.local/share/certmagic` by default) — nothing in the platform backs this up for you. Losing it means re-issuing (rate-limited, annoying) but not re-registering the ACME account (recoverable by email). Mount it on a persistent volume rather than an ephemeral container filesystem — see the volume examples in [Docker](docker.md) and [Docker Compose](docker-compose.md).

## The management API

The management API has **no in-process TLS support at all** — it only ever listens plain HTTP. There is no configuration flag that changes this. TLS for it must terminate at a reverse proxy or ingress sitting in front of it:

- [Reverse Proxy: Caddy](reverse-proxy-caddy.md) — automatic HTTPS, minimal config
- [Reverse Proxy: nginx](reverse-proxy-nginx.md) — manual certificate management (certbot or similar)

<div class="callout warning">
Whichever proxy you choose, it must set the <code>X-Forwarded-For</code> header — the management API's per-client-IP rate limit reads the client address <strong>only</strong> from that header, never from the raw TCP connection. Skipping this doesn't break anything, but it silently disables that one rate-limit dimension. See <a href="../core-concepts/rate-limits-and-quotas.html">Rate Limits and Quotas</a>.
</div>

## Next steps

- [Reverse Proxy: Caddy](reverse-proxy-caddy.md) or [Reverse Proxy: nginx](reverse-proxy-nginx.md)
- [First Boot and Upgrades](first-boot-and-upgrades.md)
