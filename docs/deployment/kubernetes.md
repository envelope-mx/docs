# Kubernetes Deployment

One Deployment + PodDisruptionBudget per role (`replicas: 2` by default — no single replica is a point of failure, and `minAvailable: 1` extends that guarantee through voluntary disruptions like node drains).

## Prerequisites

Two objects your manifests need to reference — values are deployment-specific, so create them yourself rather than expecting a default:

- **ConfigMap `envelope-config`** — non-secret settings: `ENVELOPE_DB_HOST`, `ENVELOPE_DOMAIN`, etc.
- **Secret `envelope-secrets`** — `ENVELOPE_DB_PASSWORD`, `ENVELOPE_MASTER_KEY`, `ENVELOPE_API_ADMIN_TOKEN`.

```bash
kubectl create configmap envelope-config \
  --from-literal=ENVELOPE_DB_HOST=postgres.default.svc.cluster.local \
  --from-literal=ENVELOPE_DOMAIN=mail.yourdomain.example

kubectl create secret generic envelope-secrets \
  --from-literal=ENVELOPE_DB_PASSWORD="$DB_PASSWORD" \
  --from-literal=ENVELOPE_MASTER_KEY="$(openssl rand -base64 32)" \
  --from-literal=ENVELOPE_API_ADMIN_TOKEN="$(openssl rand -base64 32)"
```

## A minimal per-role Deployment

One shape, repeated per role — this example is the `api` role; repeat with a different `ENVELOPE_ROLES` value and matching `containerPort`/probe for `smtp-inbound` (25), `smtp-submission` (587), `imap` (993), and `deliverer` (no listening port — probe `/metrics` instead, since a background-loop role has no protocol port to check).

```yaml
apiVersion: apps/v1
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
          image: ghcr.io/envelope-mx/envelope:latest # pin :vX.Y.Z in production
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
    matchLabels: { app: envelope, role: api }
```

For the protocol roles (`smtp-inbound`, `smtp-submission`, `imap`), use a TCP socket check instead of an HTTP one:

```yaml
readinessProbe:
  tcpSocket: { port: 25 }
```

```bash
kubectl apply -f envelope-api.yaml
```

## Network policy

A default-deny-ingress policy on every `app: envelope` pod, with explicit allows for Postgres reachability, the internal API port (`8081`, see below), and metrics scraping (`9090`).

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: envelope-default-deny
spec:
  podSelector: { matchLabels: { app: envelope } }
  policyTypes: [Ingress]
  ingress:
    - from: [{ podSelector: { matchLabels: { app: prometheus } } }]
      ports: [{ port: 9090 }]
```

<div class="callout gap">
Public-facing ports (25/587/993/8080) have no ingress rule in a policy like this — how traffic reaches those from outside the cluster (a Service/LoadBalancer, an ingress controller) is CNI- and operator-specific. Add an explicit allow for your actual ingress path if your CNI enforces <code>NetworkPolicy</code> on that traffic at all. If Postgres is an external/managed service (RDS, Cloud SQL, etc.) rather than in-cluster, enforce equivalent segmentation with your cloud provider's own network controls instead of a <code>NetworkPolicy</code> egress rule.
</div>

## Per-role scaling

Each role's Deployment scales independently:

```bash
kubectl scale deployment envelope-smtp-inbound --replicas=4
kubectl scale deployment envelope-deliverer --replicas=3
```

Scale `envelope-api` to scale webhook-delivery capacity along with it — dispatch runs bundled into that role, not its own Deployment.

## Optional: no live DB credential on SMTP/IMAP-facing pods

For `smtp-inbound`, `smtp-submission`, and `imap` specifically, you can eliminate live Postgres credentials from those pods' environments entirely by routing them through the `api` role's internal HTTP API instead of a direct DB connection:

- Set `ENVELOPE_INTERNAL_TOKEN_<ROLE>` (`openssl rand -base64 32` each) on the matching single-role Deployment **and** on the `api` Deployment (which needs every activated role's token to authorize incoming calls).
- Set `ENVELOPE_INTERNAL_API_URL` on the activated role, pointed at the `api` role's internal port (`8081`) — typically a ClusterIP Service in front of the `api` Deployment.

This is optional hardening on top of the base deployment above, not required to get a working instance running.

## Next steps

- [TLS and Certificates](tls-and-certificates.md)
- [Reverse Proxy: Caddy](reverse-proxy-caddy.md) or [nginx](reverse-proxy-nginx.md) — an Ingress controller typically fills this role in a cluster
- [First Boot and Upgrades](first-boot-and-upgrades.md)
