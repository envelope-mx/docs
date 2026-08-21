# Docker Compose Deployment

The simplest complete shape for a single VM: one Postgres, one rspamd, one Envelope container running every role bundled.

```yaml
# docker-compose.yaml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: envelope
      POSTGRES_PASSWORD: envelope
      POSTGRES_DB: envelope
    volumes:
      - pgdata:/var/lib/postgresql/data

  rspamd:
    image: rspamd/rspamd:latest
    ports: ["11333:11333"]

  envelope:
    image: ghcr.io/envelope-mx/envelope:latest
    depends_on: [postgres, rspamd]
    ports:
      - "25:25"
      - "587:587"
      - "993:993"
      - "8080:8080"
      - "9090:9090"
    environment:
      ENVELOPE_DB_HOST: postgres
      ENVELOPE_DB_PASSWORD: envelope
      ENVELOPE_RSPAMD_URL: http://rspamd:11333
      ENVELOPE_MASTER_KEY: ${ENVELOPE_MASTER_KEY}
      ENVELOPE_API_ADMIN_TOKEN: ${ENVELOPE_API_ADMIN_TOKEN}
      ENVELOPE_DOMAIN: mail.yourdomain.example
    volumes:
      - envelope-certs:/home/nonroot/.local/share/certmagic

volumes:
  pgdata:
  envelope-certs:
```

`:latest` tracks the newest non-prerelease tag; pin `:vX.Y.Z` (see the [releases page](https://github.com/envelope-mx/envelope/releases)) instead for anything reproducible.

## Bring it up

```bash
export ENVELOPE_MASTER_KEY=$(openssl rand -base64 32)
export ENVELOPE_API_ADMIN_TOKEN=$(openssl rand -base64 32)
docker compose up -d
curl http://localhost:8080/health
```

Store both exported values somewhere durable — see [First Boot and Upgrades](first-boot-and-upgrades.md). Losing `ENVELOPE_MASTER_KEY` after real DKIM keys and webhook secrets have been encrypted with it makes them permanently unrecoverable.

## What this shape doesn't give you

No per-role horizontal scaling, and no network segmentation between roles (the single container shares one direct database connection for everything) — fine for a solo operator or a small/medium deployment. Once you need either of those, move to [Kubernetes](kubernetes.md).

## Adding a reverse proxy to the same file

Front the management API's HTTP port with Caddy directly in this same Compose file rather than running it separately:

```yaml
  caddy:
    image: caddy:2
    depends_on: [envelope]
    ports: ["80:80", "443:443"]
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy-data:/data

volumes:
  caddy-data:
```

See [Reverse Proxy: Caddy](reverse-proxy-caddy.md) for the `Caddyfile` content this expects.

## Next steps

- [Kubernetes](kubernetes.md) — once you need to scale roles independently
- [TLS and Certificates](tls-and-certificates.md)
- [First Boot and Upgrades](first-boot-and-upgrades.md)
