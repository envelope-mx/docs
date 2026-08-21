# Deployment Overview

Envelope is distributed as a **binary and a Docker image only** — there is no public source repository to build from. Every deployment path below describes consuming one of those published artifacts.

## System shape

One codebase, split into roles selected at startup with a `--roles` flag (or your image/binary's equivalent): `api`, `smtp-inbound`, `smtp-submission`, `imap`, `deliverer`. Webhook dispatch isn't its own role — it runs as a background loop bundled into `api`, so scaling `api` replicas scales webhook delivery capacity along with it. Running with no `--roles` flag at all runs every role bundled into one process, which is the simplest shape for a single-VM or Compose deployment.

## Prerequisites

- **Postgres** — the only hard dependency. Every role needs to reach it, either directly or via the internal API (see [Kubernetes](kubernetes.md)'s network-segmentation option). Tested against Postgres 16; no version-specific features are used.
- **rspamd** (optional but recommended) — inbound's spam-scoring sidecar. Not a hard dependency: if it's unreachable, inbound mail still flows, just quarantined instead of scored.
- **Privileged ports** — `25`, `587`, `993` need root or `CAP_NET_BIND_SERVICE` to bind. Container runtimes typically grant this, or you remap ports at the Service/LoadBalancer level; for a bare-metal binary, grant the capability directly (see [Binary](binary.md)).
- **A real, registrable domain with public DNS** — needed once you move past self-signed TLS to real ACME certificates, and for inbound mail to be reachable from the internet at all (a real MX record).
- **A reverse proxy in front of the management API** — the API has **no in-process TLS support at all**. TLS for it always terminates at a proxy/ingress in front of it, never inside the process. See [TLS and Certificates](tls-and-certificates.md), [Reverse Proxy: Caddy](reverse-proxy-caddy.md), and [Reverse Proxy: nginx](reverse-proxy-nginx.md).

## Configuration reference

Everything operationally significant is an environment variable.

**Database**

| Variable | Default | Notes |
|---|---|---|
| `ENVELOPE_DB_DSN` | — | Full DSN; overrides every field below if set |
| `ENVELOPE_DB_HOST` | `127.0.0.1` | |
| `ENVELOPE_DB_PORT` | `5432` | |
| `ENVELOPE_DB_USER` / `ENVELOPE_DB_PASSWORD` | `envelope` / `envelope` | Shared credential every role falls back to |
| `ENVELOPE_DB_NAME` | `envelope` | |
| `ENVELOPE_DB_SSLMODE` | `disable` | Set `require`/`verify-full` beyond local testing |

**Required secrets**

| Variable | Notes |
|---|---|
| `ENVELOPE_MASTER_KEY` | **Required — the process refuses to boot without it.** `openssl rand -base64 32`. Encrypts DKIM private keys and webhook secrets at rest. Losing it makes every already-encrypted value permanently undecryptable — back it up as rigorously as the database itself, in a separate location. |
| `ENVELOPE_API_ADMIN_TOKEN` | Bootstrap credential authorizing every account. If unset, one is generated and logged once at boot. |

**Roles, addresses, identity**

| Variable | Default | Notes |
|---|---|---|
| `ENVELOPE_SMTP_INBOUND_ADDR` | `:25` | |
| `ENVELOPE_SMTP_SUBMISSION_ADDR` | `:587` | |
| `ENVELOPE_IMAP_ADDR` | `:993` | |
| `ENVELOPE_DOMAIN` | `localhost` | HELO/EHLO identity — set to a real domain before sending to the real internet |
| `ENVELOPE_METRICS_ADDR` | `:9090` | Every process, regardless of `--roles` |

**TLS / ACME** — see [TLS and Certificates](tls-and-certificates.md) for the full explanation.

| Variable | Notes |
|---|---|
| `ENVELOPE_ACME_EMAIL` | Unset = self-signed dev certs (default). Set = real Let's Encrypt certificates issued on demand, per active vhost. |
| `ENVELOPE_ACME_STAGING` | Set (any value) to use Let's Encrypt's staging CA while validating. |

**rspamd, quotas, rate limits**

| Variable | Default |
|---|---|
| `ENVELOPE_RSPAMD_URL` | `http://localhost:11333` |
| `ENVELOPE_RATELIMIT_IP_CAPACITY` / `_REFILL_PER_SECOND` | `20` / `1` |
| `ENVELOPE_RATELIMIT_SENDER_CAPACITY` / `_REFILL_PER_SECOND` | `10` / `0.5` |
| `ENVELOPE_RATELIMIT_API_CAPACITY` / `_REFILL_PER_SECOND` | `60` / `1` |
| `ENVELOPE_DELIVERER_PER_DOMAIN_CONCURRENCY` | `5` |

**Retention & logging**

| Variable | Default |
|---|---|
| `ENVELOPE_RETENTION_DEFAULT_DAYS` | `90` — applies to any vhost with no `retentionDays` of its own |
| `ENVELOPE_LOG_LEVEL` | `info` — every process logs structured JSON to stdout regardless |

<div class="callout note">
Every process, whatever role it runs, logs structured JSON with a shared <code>correlation_id</code> across one message's lifecycle — even as that message crosses role/replica boundaries (inbound → filter → storage → webhook, or submission → queue → deliverer → webhook). Grep for it when tracing one message's path through the system.
</div>

## Choose your path

| Path | Best for |
|---|---|
| [Binary](binary.md) | Bare VM, systemd, no container runtime |
| [Docker](docker.md) | A single container running every role |
| [Docker Compose](docker-compose.md) | One VM, Postgres + rspamd + Envelope together |
| [Kubernetes](kubernetes.md) | Per-role horizontal scaling, network segmentation |

Then [TLS and Certificates](tls-and-certificates.md), a reverse proxy ([Caddy](reverse-proxy-caddy.md) or [nginx](reverse-proxy-nginx.md)) in front of the management API, and [First Boot and Upgrades](first-boot-and-upgrades.md) to confirm everything's working.
