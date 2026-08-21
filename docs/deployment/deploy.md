# Deploy Envelope

Envelope is distributed as a **binary and a Docker image only** — there is no public source repository to build from. Every path below describes consuming one of those published artifacts.

One codebase, split into roles selected at startup with a `--roles` flag (or your image/binary's equivalent): `api`, `smtp-inbound`, `smtp-submission`, `imap`, `deliverer`. Webhook dispatch isn't its own role — it runs as a background loop bundled into `api`, so scaling `api` replicas scales webhook delivery capacity along with it. Running with no `--roles` flag at all runs every role bundled into one process, the simplest shape for a single VM or Compose deployment.

Pick a deployment method, a reverse proxy, and a TLS approach below, enter your domains, and every command on this page updates to match — copy-paste in order and you have a working instance.

## Choose your setup

<div class="wizard-panel">

<div class="wizard-panel-title">Pick your setup — the steps below update to match</div>

<div class="wizard-grid">
<div class="wizard-field">
<label for="wz-method">Deployment method</label>
<select id="wz-method">
<option value="binary">Binary (systemd)</option>
<option value="docker">Docker (single container)</option>
<option value="docker-compose" selected>Docker Compose</option>
<option value="kubernetes">Kubernetes</option>
</select>
</div>
<div class="wizard-field">
<label for="wz-proxy">Reverse proxy for the API</label>
<select id="wz-proxy">
<option value="caddy" selected>Caddy (automatic HTTPS)</option>
<option value="nginx">nginx (manual certs)</option>
<option value="none">None yet</option>
</select>
</div>
<div class="wizard-field">
<label for="wz-tls">Mail (SMTP/IMAP) TLS</label>
<select id="wz-tls">
<option value="self-signed" selected>Self-signed (dev/default)</option>
<option value="acme">Real certificates (Let's Encrypt)</option>
</select>
</div>
<div class="wizard-field">
<label for="wz-api-domain">API domain</label>
<input type="text" id="wz-api-domain" value="api.yourdomain.example" spellcheck="false" autocomplete="off">
</div>
<div class="wizard-field">
<label for="wz-mail-domain">Mail domain (ENVELOPE_DOMAIN)</label>
<input type="text" id="wz-mail-domain" value="mail.yourdomain.example" spellcheck="false" autocomplete="off">
</div>
</div>

<p class="wizard-note">Your choices and values are remembered in this browser only — nothing is sent anywhere.</p>

</div>

## Prerequisites

- **Postgres** — the only hard dependency. Every role needs to reach it, either directly or via the internal API (see the Kubernetes network-segmentation note below). Tested against Postgres 16; no version-specific features are used.
- **rspamd** (optional but recommended) — inbound's spam-scoring sidecar. Not a hard dependency: if it's unreachable, inbound mail still flows, just quarantined instead of scored.
- **Privileged ports** — `25`, `587`, `993` need root or `CAP_NET_BIND_SERVICE` to bind. <span data-when-method="binary">Grant the capability to the binary directly, or via a systemd unit's `AmbientCapabilities` (covered below).</span><span data-when-method="docker docker-compose" style="display:none">Container runtimes typically grant this automatically.</span><span data-when-method="kubernetes" style="display:none">Remap at the Service/LoadBalancer level, or grant via the container's `securityContext`.</span>
- **A real, registrable domain with public DNS** — needed once you move past self-signed TLS to real ACME certificates, and for inbound mail to be reachable from the internet at all (a real MX record).
- **A reverse proxy in front of the management API** — the API has no in-process TLS support at all. TLS for it always terminates at a proxy/ingress in front of it, never inside the process (covered below).

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

**TLS / ACME** — see [TLS for SMTP and IMAP](#tls-for-smtp-and-imap) below.

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

## Get Envelope running

<div data-when-method="binary">

<h4>1. Download the binary</h4>

<p>Grab the archive for your platform from the <a href="https://github.com/envelope-mx/envelope/releases/latest">releases page</a> and extract it. Each archive bundles the <code>envelope</code> binary alongside <code>go.mod</code> and <code>config/config.yaml</code> — keep all three together as extracted (<code>go.mod</code> is a marker file Envelope's config loader uses to find its config directory by walking up from the binary's own path, not a build input).</p>

<pre><code><span class="c-cmt"># swap linux_amd64 for linux_arm64, darwin_amd64, darwin_arm64, or windows_amd64.zip</span>
VERSION=v0.0.0  <span class="c-cmt"># see the releases page for the latest tag</span>
curl -LO "https://github.com/envelope-mx/envelope/releases/download/$VERSION/envelope_${VERSION}_linux_amd64.tar.gz"
tar -xzf "envelope_${VERSION}_linux_amd64.tar.gz"
sudo mv "envelope_${VERSION}_linux_amd64" /opt/envelope
sudo ln -sf /opt/envelope/envelope /usr/local/bin/envelope

<span class="c-cmt"># optional: verify against that release's checksums.txt</span>
curl -LO "https://github.com/envelope-mx/envelope/releases/download/$VERSION/checksums.txt"
sha256sum -c checksums.txt --ignore-missing</code></pre>

<h4>2. Grant privileged-port capability</h4>

<p>Ports 25, 587, and 993 need root or <code>CAP_NET_BIND_SERVICE</code>. Grant it to the binary itself, rather than running the whole process as root — set it on the real file in <code>/opt/envelope</code>, not the <code>/usr/local/bin</code> symlink:</p>

<pre><code>sudo setcap 'cap_net_bind_service=+ep' /opt/envelope/envelope</code></pre>

<h4>3. Configure and run</h4>

<pre><code><span class="c-cmt"># minimum required — see Configuration reference above for everything else</span>
export ENVELOPE_DB_HOST=127.0.0.1
export ENVELOPE_DB_PASSWORD=your-postgres-password
export ENVELOPE_MASTER_KEY=$(openssl rand -base64 32)
export ENVELOPE_API_ADMIN_TOKEN=$(openssl rand -base64 32)
export ENVELOPE_DOMAIN=<span class="wz" data-field="mailDomain">mail.yourdomain.example</span>

envelope --roles=api,smtp-inbound,smtp-submission,imap,deliverer
<span class="c-cmt"># omit --roles entirely to run every role in one process</span></code></pre>

<h4>4. Run as a systemd service</h4>

<p>Store <code>ENVELOPE_MASTER_KEY</code> and <code>ENVELOPE_API_ADMIN_TOKEN</code> somewhere durable before your first real boot — see the smoke-test section below.</p>

<pre><code><span class="c-cmt"># /etc/systemd/system/envelope.service</span>
[Unit]
Description=Envelope mail platform
After=network.target postgresql.service

[Service]
Type=simple
User=envelope
EnvironmentFile=/etc/envelope/envelope.env
ExecStart=/usr/local/bin/envelope
Restart=on-failure
RestartSec=5
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target</code></pre>

<p>Put every <code>ENVELOPE_*</code> variable in <code>/etc/envelope/envelope.env</code> (one <code>KEY=value</code> per line), then:</p>

<pre><code>sudo systemctl daemon-reload
sudo systemctl enable --now envelope
sudo journalctl -u envelope -f</code></pre>

<p><code>AmbientCapabilities=CAP_NET_BIND_SERVICE</code> lets the service bind privileged ports while still running as the unprivileged <code>envelope</code> user — no <code>setcap</code> needed on the binary in this case, though both approaches work together fine.</p>

<h4>5. Confirm the process itself is up</h4>

<pre><code>curl http://localhost:8080/health
<span class="c-cmt"># {"success":true,"data":{"status":"ok"}}</span></code></pre>

</div>

<div data-when-method="docker" style="display:none">

<h4>1. Pull the image</h4>

<pre><code>docker pull ghcr.io/envelope-mx/envelope:latest
<span class="c-cmt"># :latest tracks the newest non-prerelease tag; pin :vX.Y.Z for anything reproducible</span></code></pre>

<h4>2. Run it</h4>

<pre><code>docker run -d --name envelope \
  -p 25:25 -p 587:587 -p 993:993 -p 8080:8080 -p 9090:9090 \
  -e ENVELOPE_DB_HOST=host.docker.internal \
  -e ENVELOPE_DB_PASSWORD=your-postgres-password \
  -e ENVELOPE_MASTER_KEY="$(openssl rand -base64 32)" \
  -e ENVELOPE_API_ADMIN_TOKEN="$(openssl rand -base64 32)" \
  -e ENVELOPE_DOMAIN=<span class="wz" data-field="mailDomain">mail.yourdomain.example</span> \
  -v envelope-certs:/home/nonroot/.local/share/certmagic \
  ghcr.io/envelope-mx/envelope:latest</code></pre>

<p><code>host.docker.internal</code> resolves to the host machine on Docker Desktop (Mac/Windows); on native Linux, add <code>--add-host=host.docker.internal:host-gateway</code>, or point <code>ENVELOPE_DB_HOST</code> at a real reachable Postgres service instead. The image runs as a non-root distroless container (no shell, no package manager) — that's why the certmagic volume path is under <code>/home/nonroot/...</code>, not <code>/root/...</code>.</p>

<div class="callout note" data-when-tls="acme">
The <code>envelope-certs</code> volume above matters most once real ACME certificates are turned on below — without it, every container restart re-issues certificates from scratch (rate-limited, unnecessary churn). Losing the volume means re-issuing but not re-registering the ACME account (recoverable by email).
</div>

<h4>3. Confirm it's up</h4>

<pre><code>curl http://localhost:8080/health</code></pre>

<h4>Splitting roles across containers (optional)</h4>

<p>Once you outgrow a single container but aren't ready for Kubernetes, pass <code>ENVELOPE_ROLES</code> to run one role per container instead:</p>

<pre><code>docker run -d --name envelope-api      -p 8080:8080 -p 9090:9090 -e ENVELOPE_ROLES=api ...            ghcr.io/envelope-mx/envelope:latest
docker run -d --name envelope-inbound  -p 25:25     -p 9091:9090 -e ENVELOPE_ROLES=smtp-inbound ...    ghcr.io/envelope-mx/envelope:latest
docker run -d --name envelope-deliverer             -p 9092:9090 -e ENVELOPE_ROLES=deliverer ...       ghcr.io/envelope-mx/envelope:latest</code></pre>

<p>Every container still needs to reach the same Postgres and share the same <code>ENVELOPE_MASTER_KEY</code>. For most single-VM setups, switching the method above to <strong>Docker Compose</strong> is a more manageable way to express this than a series of standalone <code>docker run</code> commands.</p>

</div>

<div data-when-method="docker-compose">

<h4>docker-compose.yaml</h4>

<p>One Postgres, one rspamd (optional but recommended), one Envelope container running every role bundled. The healthcheck below makes Envelope wait for Postgres to actually accept connections rather than just for its container to start.</p>

<div data-when-tls="self-signed">

<pre><code><span class="c-cmt"># docker-compose.yaml</span>
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: envelope
      POSTGRES_PASSWORD: envelope
      POSTGRES_DB: envelope
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U envelope"]
      interval: 5s
      timeout: 5s
      retries: 10

  rspamd:
    image: rspamd/rspamd:latest
    ports: ["11333:11333"]

  envelope:
    image: ghcr.io/envelope-mx/envelope:latest
    depends_on:
      postgres:
        condition: service_healthy
      rspamd:
        condition: service_started
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
      ENVELOPE_DOMAIN: <span class="wz" data-field="mailDomain">mail.yourdomain.example</span>
    volumes:
      - envelope-certs:/home/nonroot/.local/share/certmagic

volumes:
  pgdata:
  envelope-certs:</code></pre>

</div>

<div data-when-tls="acme" style="display:none">

<pre><code><span class="c-cmt"># docker-compose.yaml — real ACME certs enabled for SMTP/IMAP</span>
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: envelope
      POSTGRES_PASSWORD: envelope
      POSTGRES_DB: envelope
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U envelope"]
      interval: 5s
      timeout: 5s
      retries: 10

  rspamd:
    image: rspamd/rspamd:latest
    ports: ["11333:11333"]

  envelope:
    image: ghcr.io/envelope-mx/envelope:latest
    depends_on:
      postgres:
        condition: service_healthy
      rspamd:
        condition: service_started
    ports:
      - "25:25"
      - "587:587"
      - "993:993"
      - "80:80"          <span class="c-cmt"># ACME HTTP-01 challenge</span>
      - "8080:8080"
      - "9090:9090"
    environment:
      ENVELOPE_DB_HOST: postgres
      ENVELOPE_DB_PASSWORD: envelope
      ENVELOPE_RSPAMD_URL: http://rspamd:11333
      ENVELOPE_MASTER_KEY: ${ENVELOPE_MASTER_KEY}
      ENVELOPE_API_ADMIN_TOKEN: ${ENVELOPE_API_ADMIN_TOKEN}
      ENVELOPE_DOMAIN: <span class="wz" data-field="mailDomain">mail.yourdomain.example</span>
      ENVELOPE_ACME_EMAIL: your-real-email@example.com
    volumes:
      - envelope-certs:/home/nonroot/.local/share/certmagic

volumes:
  pgdata:
  envelope-certs:</code></pre>

<div class="callout warning">
If Caddy/nginx are also fronting the management API on this same host, make sure they aren't already bound to port 80 — Envelope's own ACME issuance needs that port free. See <a href="#tls-for-the-management-api">TLS for the management API</a> below for how Caddy's <code>tls internal</code> mode avoids this without giving up a reverse proxy entirely.
</div>

</div>

<h4>Bring it up</h4>

<pre><code>export ENVELOPE_MASTER_KEY=$(openssl rand -base64 32)
export ENVELOPE_API_ADMIN_TOKEN=$(openssl rand -base64 32)
docker compose up -d
curl http://localhost:8080/health</code></pre>

<p>Store both exported values somewhere durable — see the smoke-test section below. Losing <code>ENVELOPE_MASTER_KEY</code> after real DKIM keys and webhook secrets have been encrypted with it makes them permanently unrecoverable.</p>

<h4>What this shape doesn't give you</h4>

<p>No per-role horizontal scaling, and no network segmentation between roles — fine for a solo operator or a small/medium deployment. Once you need either, switch the method above to <strong>Kubernetes</strong>.</p>

</div>

<div data-when-method="kubernetes" style="display:none">

<h4>Prerequisites: ConfigMap and Secret</h4>

<p>Two objects your manifests reference — values are deployment-specific, so create them yourself rather than expecting a default:</p>

<pre><code>kubectl create configmap envelope-config \
  --from-literal=ENVELOPE_DB_HOST=postgres.default.svc.cluster.local \
  --from-literal=ENVELOPE_DOMAIN=<span class="wz" data-field="mailDomain">mail.yourdomain.example</span>

kubectl create secret generic envelope-secrets \
  --from-literal=ENVELOPE_DB_PASSWORD="$DB_PASSWORD" \
  --from-literal=ENVELOPE_MASTER_KEY="$(openssl rand -base64 32)" \
  --from-literal=ENVELOPE_API_ADMIN_TOKEN="$(openssl rand -base64 32)"</code></pre>

<div data-when-tls="acme">
<p>Add <code>ENVELOPE_ACME_EMAIL</code> to the ConfigMap above too, and make sure whichever pod runs the mail-protocol roles has inbound port 80 reachable for the ACME HTTP-01 challenge.</p>
</div>

<h4>A minimal per-role Deployment</h4>

<p>One shape, repeated per role — this example is the <code>api</code> role; repeat with a different <code>ENVELOPE_ROLES</code> value and matching <code>containerPort</code>/probe for <code>smtp-inbound</code> (25), <code>smtp-submission</code> (587), <code>imap</code> (993), and <code>deliverer</code> (no listening port — probe <code>/metrics</code> instead).</p>

<pre><code>apiVersion: apps/v1
kind: Deployment
metadata:
  name: envelope-api
spec:
  replicas: 2
  selector:
    matchLabels: { app: envelope, role: api }
  template:
    metadata:
      labels: { app: envelope, role: api }
    spec:
      containers:
        - name: envelope
          image: ghcr.io/envelope-mx/envelope:latest <span class="c-cmt"># pin :vX.Y.Z in production</span>
          args: ["--roles=api"]
          envFrom:
            - configMapRef: { name: envelope-config }
            - secretRef: { name: envelope-secrets }
          ports:
            - containerPort: 8080
            - containerPort: 9090
          readinessProbe:
            httpGet: { path: /health, port: 8080 }
          livenessProbe:
            httpGet: { path: /health, port: 8080 }
---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: envelope-api
spec:
  minAvailable: 1
  selector:
    matchLabels: { app: envelope, role: api }</code></pre>

<p>For the protocol roles (<code>smtp-inbound</code>, <code>smtp-submission</code>, <code>imap</code>), use a TCP socket check instead of an HTTP one:</p>

<pre><code>readinessProbe:
  tcpSocket: { port: 25 }</code></pre>

<pre><code>kubectl apply -f envelope-api.yaml</code></pre>

<h4>Network policy</h4>

<p>A default-deny-ingress policy on every <code>app: envelope</code> pod, with explicit allows for Postgres reachability, the internal API port (8081), and metrics scraping (9090):</p>

<pre><code>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: envelope-default-deny
spec:
  podSelector: { matchLabels: { app: envelope } }
  policyTypes: [Ingress]
  ingress:
    - from: [{ podSelector: { matchLabels: { app: prometheus } } }]
      ports: [{ port: 9090 }]</code></pre>

<div class="callout gap">
Public-facing ports (25/587/993/8080) have no ingress rule in a policy like this — how traffic reaches those from outside the cluster (a Service/LoadBalancer, an ingress controller) is CNI- and operator-specific. Add an explicit allow for your actual ingress path if your CNI enforces <code>NetworkPolicy</code> on that traffic at all. If Postgres is an external/managed service (RDS, Cloud SQL, etc.) rather than in-cluster, enforce equivalent segmentation with your cloud provider's own network controls instead of a <code>NetworkPolicy</code> egress rule.
</div>

<h4>Per-role scaling</h4>

<pre><code>kubectl scale deployment envelope-smtp-inbound --replicas=4
kubectl scale deployment envelope-deliverer --replicas=3
<span class="c-cmt"># scale envelope-api to scale webhook-delivery capacity along with it — dispatch runs bundled into that role</span></code></pre>

<h4>Optional: no live DB credential on SMTP/IMAP-facing pods</h4>

<p>For <code>smtp-inbound</code>, <code>smtp-submission</code>, and <code>imap</code> specifically, eliminate live Postgres credentials from those pods entirely by routing them through the <code>api</code> role's internal HTTP API instead of a direct DB connection:</p>

<ul>
<li>Set <code>ENVELOPE_INTERNAL_TOKEN_&lt;ROLE&gt;</code> (<code>openssl rand -base64 32</code> each) on the matching single-role Deployment <strong>and</strong> on the <code>api</code> Deployment (which needs every activated role's token to authorize incoming calls).</li>
<li>Set <code>ENVELOPE_INTERNAL_API_URL</code> on the activated role, pointed at the <code>api</code> role's internal port (8081) — typically a ClusterIP Service in front of the <code>api</code> Deployment.</li>
</ul>

<p>This is optional hardening on top of the base deployment above, not required to get a working instance running.</p>

</div>

## TLS for the management API

The management API has **no in-process TLS support at all** — it only ever listens plain HTTP. There is no configuration flag that changes this. TLS for it must terminate at a reverse proxy or ingress sitting in front of it.

<div data-when-proxy="caddy">

<h4>Minimal Caddyfile</h4>

<pre><code><span class="wz" data-field="apiDomain">api.yourdomain.example</span> {
    reverse_proxy localhost:8080
}</code></pre>

<p>That's the complete configuration for automatic HTTPS plus a correctly-forwarded client IP — Caddy sets <code>X-Forwarded-For</code> by default, which is what the <a href="../core-concepts/rate-limits-and-quotas.html">per-IP API rate limit</a> reads. Caddy obtains and renews the certificate on its own via Let's Encrypt.</p>

<pre><code><span class="c-cmt"># as a system service, alongside a Binary deployment</span>
sudo apt install caddy
sudo cp Caddyfile /etc/caddy/Caddyfile
sudo systemctl restart caddy</code></pre>

<p>Alongside Docker Compose, add a <code>caddy</code> service to the same file instead:</p>

<pre><code>  caddy:
    image: caddy:2
    depends_on: [envelope]
    ports: ["80:80", "443:443"]
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy-data:/data

volumes:
  caddy-data:</code></pre>

<div class="callout warning">
If Envelope's own ACME issuance is also turned on above on this same host, Caddy's <code>80:80</code> here will conflict with it — both want that port for their own certificate challenges. Either run Caddy with <code>tls internal</code> instead (its own self-managed local CA, no port 80/443 needed at all — trades the public trust chain for a port that's actually free) on a different external port, or put the two on separate hosts.
</div>

<h4>Confirm it's working</h4>

<pre><code>curl https://<span class="wz" data-field="apiDomain">api.yourdomain.example</span>/health</code></pre>

<p>A trusted-certificate <code>200</code> means Caddy is correctly terminating TLS and forwarding to the API.</p>

</div>

<div data-when-proxy="nginx" style="display:none">

<h4>Minimal server block</h4>

<pre><code>server {
    listen 443 ssl;
    server_name <span class="wz" data-field="apiDomain">api.yourdomain.example</span>;

    ssl_certificate     /etc/letsencrypt/live/<span class="wz" data-field="apiDomain">api.yourdomain.example</span>/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/<span class="wz" data-field="apiDomain">api.yourdomain.example</span>/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 80;
    server_name <span class="wz" data-field="apiDomain">api.yourdomain.example</span>;
    return 301 https://$host$request_uri;
}</code></pre>

<div class="callout warning">
The <code>proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;</code> line is not optional — Envelope's per-client-IP API rate limit reads the client address <strong>only</strong> from that header. Omitting it doesn't break requests, it silently disables that rate-limit dimension for every caller. See <a href="../core-concepts/rate-limits-and-quotas.html">Rate Limits and Quotas</a>.
</div>

<h4>Obtain the certificate with certbot</h4>

<pre><code>sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d <span class="wz" data-field="apiDomain">api.yourdomain.example</span>
<span class="c-cmt"># certbot offers to edit the server block above in place and sets up automatic renewal</span></code></pre>

<pre><code>sudo nginx -t && sudo systemctl reload nginx  <span class="c-cmt"># reload after any config change</span></code></pre>

<div class="callout warning">
Unlike Caddy, nginx's port 80/443 usage above will also conflict with Envelope's own ACME issuance (above) if both run on the same host — certbot's port-80 renewal challenge and Envelope's HTTP-01 challenge both want it. Put them on separate hosts, or keep mail TLS self-signed on this host.
</div>

<h4>Confirm it's working</h4>

<pre><code>curl https://<span class="wz" data-field="apiDomain">api.yourdomain.example</span>/health</code></pre>

</div>

<div data-when-proxy="none" style="display:none">

<div class="callout warning">
The management API has no in-process TLS support at all — this isn't optional. Anything talking to it in production should go through a reverse proxy or ingress, not directly to port 8080. Switch the proxy choice above to see a minimal Caddy or nginx configuration once you're ready.
</div>

<h4>Confirm the API itself responds (unencrypted, local only)</h4>

<pre><code>curl http://localhost:8080/health</code></pre>

</div>

<div data-when-method="kubernetes">

<div class="callout note">
In a cluster, an Ingress controller typically fills the reverse-proxy role above instead of a standalone Caddy/nginx process — adapt the server block/Caddyfile shown for your current proxy choice into Ingress annotations (or an equivalent <code>cert-manager</code> <code>Issuer</code>) rather than running Caddy/nginx as its own Deployment.
</div>

</div>

## TLS for SMTP and IMAP

Envelope's TLS story has two completely separate halves — the section above covers the management API; this one covers the mail protocols, which handle their own TLS in-process.

<div data-when-tls="self-signed">

<p><strong>Default behavior</strong> — a self-signed certificate, generated per-process. Fine for local testing and internal use; browsers and mail clients won't trust it, and neither will remote MTAs verifying inbound TLS on delivery. No environment variables needed for this.</p>

</div>

<div data-when-tls="acme" style="display:none">

<p>Set <code>ENVELOPE_ACME_EMAIL</code> (and a real <code>ENVELOPE_DOMAIN</code> — <span class="wz" data-field="mailDomain">mail.yourdomain.example</span> above) on every <code>smtp-inbound</code>/<code>smtp-submission</code>/<code>imap</code> process. Certificates are then issued automatically on first handshake, <strong>per active vhost domain specifically</strong> — issuance is scoped to domains actually registered as active vhosts, which is what prevents anyone from burning your ACME rate-limit budget requesting certificates for arbitrary hostnames. Needs port 80 reachable from the public internet (the ACME HTTP-01 challenge).</p>

<pre><code>ENVELOPE_ACME_EMAIL=your-real-email@example.com
<span class="c-cmt"># validate first — exercises the real issuance flow against Let's Encrypt's staging CA</span>
<span class="c-cmt"># (untrusted certs, doesn't spend production rate-limit budget); remove once confirmed:</span>
ENVELOPE_ACME_STAGING=1</code></pre>

<div class="callout warning">
Setting <code>ENVELOPE_ACME_EMAIL</code> is your explicit agreement to the CA's subscriber agreement — there's no separate prompt. Certificate/account state persists to local storage inside the container/host (<code>~/.local/share/certmagic</code> by default, <code>/home/nonroot/.local/share/certmagic</code> in the Docker image) — nothing in the platform backs this up for you. Mount it on a persistent volume rather than an ephemeral container filesystem (see the Docker/Compose steps above). Losing it means re-issuing (rate-limited, annoying) but not re-registering the ACME account (recoverable by email).
</div>

</div>

## First boot and smoke test

<p>Schema migrations run automatically and unconditionally on every boot of a process holding a direct database connection — idempotent, so this is safe even when multiple replicas start concurrently. There is no separate migration command to invoke yourself, on first boot or any subsequent deploy.</p>

<pre><code><span class="c-cmt"># 1. Liveness</span>
curl https://<span class="wz" data-field="apiDomain">api.yourdomain.example</span>/health

<span class="c-cmt"># 2. Create a test account (auto-issues its first token)</span>
curl -X POST https://<span class="wz" data-field="apiDomain">api.yourdomain.example</span>/accounts \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"name": "Smoke Test"}'
<span class="c-cmt"># → save data.token.token as $ACCOUNT_TOKEN, data.account.id as $ACCOUNT_ID</span>

<span class="c-cmt"># 3. Self-serve create a vhost (no admin token needed from here on)</span>
curl -X POST https://<span class="wz" data-field="apiDomain">api.yourdomain.example</span>/accounts/$ACCOUNT_ID/vhosts \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"domain": "<span class="wz" data-field="mailDomain">mail.yourdomain.example</span>"}'

<span class="c-cmt"># 4. Send via the REST API (no mailbox needed)</span>
curl -X POST https://<span class="wz" data-field="apiDomain">api.yourdomain.example</span>/accounts/$ACCOUNT_ID/messages \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"from": "test@<span class="wz" data-field="mailDomain">mail.yourdomain.example</span>", "to": ["you@example.com"],
       "subject": "Envelope smoke test", "text": "If you got this, the deployment works."}'

<span class="c-cmt"># 5. Confirm metrics are being scraped</span>
curl http://localhost:9090/metrics | head</code></pre>

<p>This exercises the full self-serve onboarding sequence, the outbound send pipeline, and DKIM signing in one pass. See the <a href="../getting-started/quickstart.html">Quickstart</a> for the same flow with fuller response examples, and <a href="../guides/set-up-a-new-vhost.html">Set up a new vhost</a> for a version that goes on to configure policy and DNS for real production sending.</p>

## Upgrading

<div data-when-method="binary">

<ol>
<li>Download the new release (repeat step 1 above).</li>
<li>Restart the service — migrations run automatically as it boots: <code>sudo systemctl restart envelope</code></li>
</ol>

</div>

<div data-when-method="docker" style="display:none">

<pre><code>docker pull ghcr.io/envelope-mx/envelope:latest
docker stop envelope && docker rm envelope
<span class="c-cmt"># re-run the same `docker run` command from step 2 above</span></code></pre>

</div>

<div data-when-method="docker-compose">

<pre><code>docker compose pull envelope
docker compose up -d</code></pre>

</div>

<div data-when-method="kubernetes" style="display:none">

<p>A standard rolling update (the default <code>Deployment</code> strategy) works: new-version pods come up, pass their readiness probe, and start serving before old-version pods terminate — the automatic per-boot migration is safe under this overlap since it's idempotent.</p>

<pre><code>kubectl set image deployment/envelope-api envelope=ghcr.io/envelope-mx/envelope:v1.2.3
<span class="c-cmt"># repeat per role Deployment</span></code></pre>

</div>

<p>Whatever the method, watch <code>/metrics</code> and your logs for the first few minutes after each role's rollout — a schema or config incompatibility in a new release would typically surface as boot failures or a spike in <code>5xx</code>s immediately, not as a silent delayed failure.</p>

<div class="callout note">
There is no documented downgrade path. If a new release adds a schema migration, rolling back to the previous binary/image against an already-migrated database is not something this platform guarantees works. Keep a recent database backup before upgrading a production deployment.
</div>

## Next steps

- [Multi-Tenancy](../core-concepts/multi-tenancy.md) — what to build on top of a freshly-booted instance
- [Guides](../guides/send-via-rest.md) — worked end-to-end integration walkthroughs
- [Rate Limits and Quotas](../core-concepts/rate-limits-and-quotas.md) — the operator-side environment variables behind the IP-based limits
