# Authentication

Every management API request requires a bearer token: `Authorization: Bearer <token>`. There are exactly two kinds of credential.

## Admin token

The platform bootstrap credential, sourced from the operator's `ENVELOPE_API_ADMIN_TOKEN` environment variable. It authorizes **every** account. If the operator didn't set it, the process generates one and logs it once at boot — fine for throwaway testing, not something to rely on surviving a restart.

The admin token is used for exactly two things in normal operation:

- `POST /accounts` — create a new tenant account (the only admin-gated step in onboarding)
- `GET /accounts` and `GET /vhosts` — list every account/vhost across the whole deployment

Everything else an admin token can do, an account's own token can also do for its own account — the admin token is a superset, not a separate API surface.

## Account tokens

Scoped to one account via `Token.accountId`. A token minted for an account can manage **every vhost that account owns** — there is currently no narrower per-vhost token scope. Tokens are formatted `env_<43-character base64url string>` (32 random bytes) and stored only as a SHA-256 hash — the raw value is shown exactly once, at creation, and can never be retrieved again.

Mint further tokens, list them, or revoke one via [API Reference → Tokens](../api-reference/tokens.md). Revoking one token has no effect on any other token minted for the same account.

## The self-serve onboarding sequence

This is the complete flow from zero to a sending vhost. Only step 1 needs the admin token.

1. **(Admin, one-time)** `POST /accounts {"name": "..."}` → response includes the new account **and** its auto-issued first token (label `"default"`).
2. **(Self-serve)** `POST /accounts/:accountId/vhosts {"domain": "..."}` → creates the vhost and generates its DKIM key pair in the same step.
3. Publish the DKIM DNS record (see [DNS and DKIM](../core-concepts/dns-and-dkim.md)).
4. **(Optional)** `POST /vhosts/:vhostId/mailboxes` — only needed if you plan to use SMTP submission or IMAP; REST send and webhooks don't require a mailbox at all.
5. **(Optional)** `POST /accounts/:accountId/tokens {"label": "..."}` — mint additional tokens, e.g. one per service or environment. Every token minted this way is scoped to the whole account.
6. **(Optional)** `PATCH /vhosts/:id/policy` — set real `maxMessageBytes`, `dailyQuota`, spam thresholds, and `retentionDays`; all default to `0` (unconfigured) otherwise.
7. **(Optional)** `POST /vhosts/:vhostId/webhooks` — subscribe to delivery/receipt events.

<div class="callout gap">
There is currently no self-serve account creation — creating the first account for a new tenant is always an admin-only, one-time operation. There is also no per-vhost token scope: any token minted for an account can act on every vhost that account owns.
</div>

## Authorization rules

Every authenticated route falls into one of three checks:

- **Admin-only** — `POST`/`GET /accounts`, `GET /vhosts` (the full cross-tenant list).
- **Account-scoped** — admin, or a token whose account matches the `:accountId` in the URL. Used for `GET /accounts/:id`, the account audit log, self-serve vhost creation, token CRUD, and sending a message.
- **Vhost-scoped (via account ownership)** — admin, or a token belonging to the account that owns the `:vhostId` in the URL. Used for every route nested under `/vhosts/:vhostId/...` (mailboxes, webhooks, the vhost audit log, data export/delete, deactivate, policy update).

<div class="callout note">
A non-admin token probing a vhost ID that doesn't belong to its account gets the same <code>403</code> whether that ID belongs to a different tenant or doesn't exist at all — never a distinguishable <code>404</code>. This is deliberate: a <code>404</code>/<code>403</code> split would let any authenticated account enumerate which vhost IDs exist across every other tenant. See <a href="../core-concepts/errors-and-responses.html">Errors and Responses</a>.
</div>

## Next steps

- [Multi-Tenancy](../core-concepts/multi-tenancy.md) — the account/vhost/mailbox ownership model in full
- [API Reference → Accounts](../api-reference/accounts.md) and [→ Tokens](../api-reference/tokens.md)
- [Rotate a token](../guides/rotate-a-token.md) — a worked zero-downtime rotation walkthrough
