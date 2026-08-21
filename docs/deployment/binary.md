# Binary Deployment

Running the published `envelope` binary directly on a VM or bare metal — no container runtime involved.

## 1. Download

Download the release for your platform from the releases page (`https://github.com/<org>/envelope/releases`) and place it on your `PATH`:

```bash
curl -Lo envelope https://github.com/<org>/envelope/releases/latest/download/envelope-linux-amd64
chmod +x envelope
sudo mv envelope /usr/local/bin/envelope
```

## 2. Grant privileged-port capability

Ports 25, 587, and 993 need root or `CAP_NET_BIND_SERVICE`. Rather than running the whole process as root, grant the capability to the binary itself:

```bash
sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/envelope
```

## 3. Configure

Every setting is an environment variable — see the [full reference](overview.md#configuration-reference). At minimum:

```bash
export ENVELOPE_DB_HOST=127.0.0.1
export ENVELOPE_DB_PASSWORD=your-postgres-password
export ENVELOPE_MASTER_KEY=$(openssl rand -base64 32)
export ENVELOPE_API_ADMIN_TOKEN=$(openssl rand -base64 32)
export ENVELOPE_DOMAIN=mail.yourdomain.example
```

Store `ENVELOPE_MASTER_KEY` and `ENVELOPE_API_ADMIN_TOKEN` somewhere durable before your first real boot — see [First Boot and Upgrades](first-boot-and-upgrades.md).

## 4. Run

```bash
envelope --roles=api,smtp-inbound,smtp-submission,imap,deliverer
```

Omitting `--roles` entirely runs every role in one process — the simplest shape for a single machine. Split roles across multiple machines/processes by running the binary once per role, each pointed at the same Postgres.

## 5. Run as a systemd service

```ini
# /etc/systemd/system/envelope.service
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
WantedBy=multi-user.target
```

Put every `ENVELOPE_*` variable in `/etc/envelope/envelope.env` (one `KEY=value` per line, not committed anywhere with secrets in it), then:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now envelope
sudo journalctl -u envelope -f
```

`AmbientCapabilities=CAP_NET_BIND_SERVICE` lets the service bind privileged ports while still running as the unprivileged `envelope` user — no `setcap` needed on the binary in this case, though both approaches work.

## 6. Confirm it's up

```bash
curl http://localhost:8080/health
# {"success":true,"data":{"status":"ok"}}
```

## Next steps

- [TLS and Certificates](tls-and-certificates.md)
- [Reverse Proxy: Caddy](reverse-proxy-caddy.md) or [nginx](reverse-proxy-nginx.md) — required in front of the management API
- [First Boot and Upgrades](first-boot-and-upgrades.md)
