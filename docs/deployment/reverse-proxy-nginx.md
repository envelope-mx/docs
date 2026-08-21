# Reverse Proxy: nginx

The management API has no in-process TLS (see [TLS and Certificates](tls-and-certificates.md)). Unlike Caddy, nginx doesn't obtain certificates automatically — pair it with certbot (or your own certificate management) for the TLS material, and make sure to forward the client IP explicitly.

## Minimal server block

```nginx
server {
    listen 443 ssl;
    server_name api.yourdomain.example;

    ssl_certificate     /etc/letsencrypt/live/api.yourdomain.example/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.example/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 80;
    server_name api.yourdomain.example;
    return 301 https://$host$request_uri;
}
```

<div class="callout warning">
The <code>proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;</code> line is not optional — Envelope's per-client-IP API rate limit reads the client address <strong>only</strong> from that header, never from the raw TCP connection nginx sees. Omitting it doesn't break requests, it silently disables that one rate-limit dimension for every caller. See <a href="../core-concepts/rate-limits-and-quotas.html">Rate Limits and Quotas</a>.
</div>

## Obtaining the certificate with certbot

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d api.yourdomain.example
```

certbot will offer to edit the server block above in place and set up automatic renewal via a systemd timer or cron entry — accept that unless you're managing certificates another way.

## Reload after any config change

```bash
sudo nginx -t && sudo systemctl reload nginx
```

## Confirm it's working

```bash
curl https://api.yourdomain.example/health
```

<div class="callout note">
This config only fronts the management API's HTTP port (8080). SMTP (25/587) and IMAP (993) are raw TCP protocols with their own independent TLS story — see <a href="tls-and-certificates.html">TLS and Certificates</a> — Envelope's built-in ACME support handles those directly. nginx's <code>stream {}</code> block can proxy raw TCP if you want a single centralized ingress point, but that's an advanced setup outside what Envelope expects by default.
</div>

## Next steps

- [First Boot and Upgrades](first-boot-and-upgrades.md)
- [Reverse Proxy: Caddy](reverse-proxy-caddy.md) — the lower-effort alternative if you don't already run nginx
