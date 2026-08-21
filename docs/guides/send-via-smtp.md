# Guide: Send via SMTP

Using [authenticated SMTP submission](../smtp/submission.md) instead of REST — useful when you already have mail-sending code that speaks SMTP and don't want to rewrite it against a new API.

## Prerequisites

A [mailbox](../api-reference/mailboxes.md) on an active vhost — SMTP submission authenticates against a mailbox credential, not an account token.

```bash
curl -X POST https://mail.yourdomain.example/vhosts/$VHOST_ID/mailboxes \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"localPart": "billing", "password": "a strong password"}'
```

## Sending with Python's `smtplib`

```python
import smtplib
import ssl
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

msg = MIMEMultipart("alternative")
msg["From"] = "billing@mail.acme.example"
msg["To"] = "customer@example.com"
msg["Subject"] = "Your invoice for August"
msg.attach(MIMEText("Plain-text body.", "plain"))
msg.attach(MIMEText("<p>HTML body.</p>", "html"))

context = ssl.create_default_context()
with smtplib.SMTP("mail.yourdomain.example", 587) as server:
    server.starttls(context=context)
    server.login("billing@mail.acme.example", "a strong password")
    server.sendmail(msg["From"], [msg["To"]], msg.as_string())
```

`login`'s username is the **full mailbox address**, not just the local part — this is how Envelope knows which vhost's DKIM key to sign the outgoing message with.

## Sending with nodemailer

```javascript
const nodemailer = require("nodemailer");

const transporter = nodemailer.createTransport({
  host: "mail.yourdomain.example",
  port: 587,
  secure: false, // STARTTLS, not implicit TLS
  auth: { user: "billing@mail.acme.example", pass: "a strong password" },
});

await transporter.sendMail({
  from: "billing@mail.acme.example",
  to: "customer@example.com",
  subject: "Your invoice for August",
  text: "Plain-text body.",
  html: "<p>HTML body.</p>",
});
```

## Common mistakes

- **Sending before `AUTH` succeeds** → `530 5.7.0 authentication required`. Confirm your library actually calls `STARTTLS` then `AUTH` before `MAIL FROM` — most do this automatically, but some need it configured explicitly.
- **Using `LOGIN`-only auth libraries** — Envelope's submission server only advertises `PLAIN`. A client hardcoded to only try `AUTH LOGIN` will fail to authenticate at all.
- **`From:` not matching the authenticated mailbox's vhost** — a mailbox can only send DKIM-signed as its own vhost's domain, regardless of what the message's `From:` header claims.

## Next steps

- [Send via REST](send-via-rest.md) — the alternative if you'd rather not manage SMTP client configuration
- [Receive and verify webhooks](receive-and-verify-webhooks.md)
