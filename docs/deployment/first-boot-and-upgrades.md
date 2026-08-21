# First Boot and Upgrades

## Migrations run automatically

Schema migrations run automatically and unconditionally on every boot of a process holding a direct database connection — idempotent, so this is safe even when multiple replicas start concurrently and each runs it. There is no separate migration command to invoke yourself, on first boot or any subsequent deploy.

## Smoke-testing a fresh deployment

```bash
# 1. Liveness
curl https://api.yourdomain.example/health

# 2. Create a test account (auto-issues its first token)
curl -X POST https://api.yourdomain.example/accounts \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"name": "Smoke Test"}'
# → save data.token.token as $ACCOUNT_TOKEN, data.account.id as $ACCOUNT_ID

# 3. Self-serve create a vhost (no admin token needed from here on)
curl -X POST https://api.yourdomain.example/accounts/$ACCOUNT_ID/vhosts \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"domain": "smoke-test.yourdomain.example"}'

# 4. Send via the REST API (no mailbox needed)
curl -X POST https://api.yourdomain.example/accounts/$ACCOUNT_ID/messages \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"from": "test@smoke-test.yourdomain.example", "to": ["you@example.com"],
       "subject": "Envelope smoke test", "text": "If you got this, the deployment works."}'

# 5. Confirm metrics are being scraped
curl http://localhost:9090/metrics | head
```

This exercises the full self-serve onboarding sequence, the outbound send pipeline, and DKIM signing in one pass. See the [Quickstart](../getting-started/quickstart.md) for the same flow with fuller response examples, and [Set up a new vhost](../guides/set-up-a-new-vhost.md) for a version that goes on to configure policy and DNS for real production sending.

## Upgrading

1. Pull the new image tag / download the new binary release.
2. Restart each role's process/pod with the new version — migrations run automatically as each one boots, so there's no separate migration step to sequence beforehand.
3. In Kubernetes, a standard rolling update (the default `Deployment` strategy) works: new-version pods come up, pass their readiness probe, and start serving before old-version pods terminate — the automatic per-boot migration is safe under this overlap since it's idempotent.
4. Watch `/metrics` and your logs for the first few minutes after each role's rollout — a schema or config incompatibility in a new release would typically surface as boot failures or a spike in `5xx`s immediately, not as a silent delayed failure.

<div class="callout note">
There is no documented downgrade path. If a new release adds a schema migration, rolling back to the previous binary/image against an already-migrated database is not something this platform guarantees works. Keep a recent database backup before upgrading a production deployment.
</div>

## Next steps

- [Multi-Tenancy](../core-concepts/multi-tenancy.md) — what to build on top of a freshly-booted instance
- [Guides](../guides/send-via-rest.md) — worked end-to-end integration walkthroughs
