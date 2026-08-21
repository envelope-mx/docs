# Binary Deployment

Running the published `envelope` binary directly on a VM or bare metal — no container runtime involved.

## 1. Download

Grab the archive for your platform from the [releases page](https://github.com/envelope-mx/envelope/releases/latest) and extract it. Each archive bundles the `envelope` binary alongside `go.mod` and `config/config.yaml` — keep all three together as extracted (`go.mod` is a marker file, not a build input: Envelope's Goose-based config loader finds its config directory by walking up from the binary's own path looking for `go.mod`, so it can't be moved to `/usr/local/bin` on its own without leaving that marker behind):

```bash
VERSION=v0.0.0  # see the releases page for the latest tag
curl -LO "https://github.com/envelope-mx/envelope/releases/download/$VERSION/envelope_${VERSION}_linux_amd64.tar.gz"
tar -xzf "envelope_${VERSION}_linux_amd64.tar.gz"
sudo mv "envelope_${VERSION}_linux_amd64" /opt/envelope
sudo ln -sf /opt/envelope/envelope /usr/local/bin/envelope
```

Swap `linux_amd64` for `linux_arm64`, `darwin_amd64`, or `darwin_arm64` for other platforms (Windows ships as a `.zip` with the same layout). Optionally verify the download against that release's `checksums.txt`:

```bash
curl -LO "https://github.com/envelope-mx/envelope/releases/download/$VERSION/checksums.txt"
sha256sum -c checksums.txt --ignore-missing
```

## 2. Grant privileged-port capability

Ports 25, 587, and 993 need root or `CAP_NET_BIND_SERVICE`. Rather than running the whole process as root, grant the capability to the binary itself — set it on the real file in `/opt/envelope`, not the `/usr/local/bin` symlink:

```bash
sudo setcap 'cap_net_bind_service=+ep' /opt/envelope/envelope
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
