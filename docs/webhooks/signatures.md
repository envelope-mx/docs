# Signatures

Every webhook delivery carries `X-Envelope-Signature: sha256=<hex>` — an HMAC-SHA256 of the **raw request body**, keyed with the subscription's secret (the one you supplied when creating the subscription — it's write-only, never returned by the API afterward, so verification requires you to have kept your own copy).

## Verifying in shell (illustrative)

```bash
body='{"id":"evt_...","type":"message.delivered","vhost":"mail.acme.example","createdAt":"...","data":{}}'
secret='your-webhook-secret'

expected="sha256=$(printf '%s' "$body" | openssl dgst -sha256 -hmac "$secret" | sed 's/^.* //')"
received='sha256=...'  # from the X-Envelope-Signature header

[ "$expected" = "$received" ] && echo "valid" || echo "INVALID — reject this request"
```

## What to get right

- **Compare the raw bytes, not a re-serialized/re-parsed body.** Compute the HMAC over exactly what arrived on the wire, before your framework parses it into an object — re-encoding JSON can reorder keys or change whitespace and silently break verification even for a genuine delivery.
- **Use a constant-time comparison**, not `==` / `strcmp`, to avoid a timing side-channel on the comparison itself. Most languages have one built in (`hmac.compare_digest` in Python, `crypto.timingSafeEqual` in Node, `subtle.ConstantTimeCompare` in Go).
- **Reject on mismatch with a non-`2xx` status.** A failed verification should look like a failed delivery to Envelope's retry logic, not a silent no-op.
- **Deduplicate by `id`, not by arrival.** A delivery that timed out on your end but actually succeeded, or one being manually [redriven](../api-reference/webhooks.md), can validly arrive more than once with the identical `id` and payload — treat processing the same `id` twice as a no-op rather than double-applying its effect.

## Next steps

- [Events](events.md)
- [Receive and verify webhooks](../guides/receive-and-verify-webhooks.md) — the same logic wired into a small real endpoint
