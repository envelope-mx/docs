# Status Codes

## Management API (HTTP)

| Code | Meaning | Where it appears |
|---|---|---|
| <span class="method-badge get">200</span> | OK | Successful `GET`, and mutations that don't create a resource (`PATCH`/`DELETE`) |
| <span class="method-badge post">201</span> | Created | Every `POST` that creates a resource |
| <span class="method-badge patch">400</span> | Bad Request | Missing/malformed field, invalid attachment base64, reserved-header collision |
| <span class="method-badge patch">401</span> | Unauthorized | Missing/malformed `Authorization` header, unrecognized or revoked token |
| <span class="method-badge delete">403</span> | Forbidden | Admin required but caller isn't; token doesn't own the target account/vhost |
| <span class="method-badge delete">404</span> | Not Found | Resource genuinely missing (admin caller, or an account-level lookup) |
| <span class="method-badge patch">429</span> | Too Many Requests | API per-IP rate limit; REST send daily quota exceeded |
| <span class="method-badge delete">500</span> | Internal Server Error | Any unhandled failure, including duplicate-domain/duplicate-mailbox conflicts (see [Known Limitations](known-limitations.md)) |
| <span class="method-badge delete">503</span> | Service Unavailable | Webhook redrive with no dispatcher configured; REST send to a vhost with no DKIM key |

Full envelope shape and pagination convention: [Errors and Responses](../core-concepts/errors-and-responses.md).

## SMTP

| Code | Meaning | Where it appears |
|---|---|---|
| `421` | Service busy / try again | Transport-level temporary failure a remote MTA may return to the deliverer |
| `450 4.7.1` | Temporary — rate limited | Inbound per-IP/per-sender rate limit |
| `452 4.7.1` | Temporary — quota exceeded | Submission daily sending quota exceeded |
| `451` | Temporary — internal failure | No DKIM key configured; DKIM signing failed; storage/queue failure |
| `501 5.1.3` | Malformed recipient address | Inbound `RCPT TO` with no `@` |
| `530 5.7.0` | Authentication required | Submission command attempted before `AUTH` succeeded |
| `550 5.1.1` | No such domain / user unknown | Inbound `RCPT TO` for a non-vhost domain; a real bounce reason from a remote MTA |
| `550 5.7.1` | Message rejected | Inbound filter pipeline reject verdict (spam score or DMARC `p=reject`) |
| `552 5.3.4` | Message too large | Inbound `DATA` exceeding the effective size ceiling |

See [SMTP: Inbound](../smtp/inbound.md) and [SMTP: Submission](../smtp/submission.md) for the full per-command behavior behind each of these.

## Next steps

- [Errors and Responses](../core-concepts/errors-and-responses.md)
- [Known Limitations](known-limitations.md)
