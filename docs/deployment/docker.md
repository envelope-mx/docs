# Docker Deployment

Running the published Docker image directly, with `docker run` — one container bundling every role.

## Pull the image

```bash
docker pull ghcr.io/envelope-mx/envelope:latest
```

`:latest` tracks the newest non-prerelease tag; pin `:vX.Y.Z` (see the [releases page](https://github.com/envelope-mx/envelope/releases)) instead for anything reproducible.

## Run

```bash
docker run -d --name envelope \
  -p 25:25 -p 587:587 -p 993:993 -p 8080:8080 -p 9090:9090 \
  -e ENVELOPE_DB_HOST=host.docker.internal \
  -e ENVELOPE_DB_PASSWORD=your-postgres-password \
  -e ENVELOPE_MASTER_KEY="$(openssl rand -base64 32)" \
  -e ENVELOPE_API_ADMIN_TOKEN="$(openssl rand -base64 32)" \
  -e ENVELOPE_DOMAIN=mail.yourdomain.example \
  -v envelope-certs:/home/nonroot/.local/share/certmagic \
  ghcr.io/envelope-mx/envelope:latest
```

`host.docker.internal` resolves to the host machine on Docker Desktop (Mac/Windows); on native Linux, add `--add-host=host.docker.internal:host-gateway` to the run command, or point `ENVELOPE_DB_HOST` at a real reachable Postgres service instead.

The image runs as a non-root distroless container (no shell, no package manager) — this is why the certmagic volume path above is under `/home/nonroot/...` rather than `/root/...`.

## Persist ACME certificate state

If you set `ENVELOPE_ACME_EMAIL` for real Let's Encrypt certificates (see [TLS and Certificates](tls-and-certificates.md)), mount a volume at the certmagic storage path (shown above) — otherwise every container restart re-issues certificates from scratch, which is rate-limited and unnecessary churn. Losing this volume means re-issuing (annoying, rate-limited) but not re-registering the ACME account (recoverable by email).

## Confirm it's up

```bash
curl http://localhost:8080/health
```

## Splitting roles across containers

Pass `--roles=<role>` (or your image's equivalent env var/entrypoint arg) to run one role per container instead of everything bundled — useful once you outgrow a single-container deployment but aren't ready for [Kubernetes](kubernetes.md):

```bash
docker run -d --name envelope-api    -p 8080:8080 -p 9090:9090 -e ENVELOPE_ROLES=api ...            ghcr.io/envelope-mx/envelope:latest
docker run -d --name envelope-inbound -p 25:25    -p 9091:9090 -e ENVELOPE_ROLES=smtp-inbound ...    ghcr.io/envelope-mx/envelope:latest
docker run -d --name envelope-deliverer            -p 9092:9090 -e ENVELOPE_ROLES=deliverer ...      ghcr.io/envelope-mx/envelope:latest
```

Every container still needs to reach the same Postgres and share the same `ENVELOPE_MASTER_KEY`. For most single-VM deployments, [Docker Compose](docker-compose.md) is a more manageable way to express this than a series of standalone `docker run` commands.

## Next steps

- [Docker Compose](docker-compose.md) — Postgres, rspamd, and Envelope wired together declaratively
- [TLS and Certificates](tls-and-certificates.md)
- [Reverse Proxy: Caddy](reverse-proxy-caddy.md) or [nginx](reverse-proxy-nginx.md)
