package builder

const sharedHead = `<link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">
    <link rel="icon" href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Crect width='24' height='24' rx='6' fill='%23f2b544'/%3E%3Cpath d='M4 6a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2z' fill='none' stroke='%230a0a0a' stroke-width='2'/%3E%3Cpolyline points='4 7 12 13 20 7' fill='none' stroke='%230a0a0a' stroke-width='2'/%3E%3C/svg%3E">`

const themeSwitchHTML = `<div class="theme-switch" role="group" aria-label="Theme">
                <button class="theme-switch-btn" data-theme-choice="system" aria-label="Match system theme" title="System">
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>
                </button>
                <button class="theme-switch-btn" data-theme-choice="light" aria-label="Light theme" title="Light">
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"></circle><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"></path></svg>
                </button>
                <button class="theme-switch-btn" data-theme-choice="dark" aria-label="Dark theme" title="Dark">
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
                </button>
            </div>`

const sidebarFooterHTML = `<div class="sidebar-footer">
            <a href="https://github.com/envelope-mx/docs/issues" target="_blank" rel="noopener" class="sidebar-footer-link">
                <span>Support</span>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M7 17L17 7"></path><path d="M7 7h10v10"></path></svg>
            </a>
            <a href="https://github.com/envelope-mx/docs" target="_blank" rel="noopener" class="sidebar-footer-link">
                <span>GitHub</span>
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M7 17L17 7"></path><path d="M7 7h10v10"></path></svg>
            </a>
        </div>`

const searchBlockHTML = `<div class="search-container">
            <div class="search-input-wrapper">
                <svg class="search-icon" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="11" cy="11" r="8"></circle>
                    <path d="M21 21l-4.35-4.35"></path>
                </svg>
                <input type="text" class="search-input" placeholder="Search..." aria-label="Search documentation">
                <kbd class="search-shortcut">⌘K</kbd>
            </div>
            <div class="search-results"></div>
        </div>`

const searchModalHTML = `<div class="search-modal" role="dialog" aria-modal="true" aria-label="Search documentation">
        <div class="search-modal-content">
            <div class="search-modal-header">
                <svg class="search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="11" cy="11" r="8"></circle>
                    <path d="M21 21l-4.35-4.35"></path>
                </svg>
                <input type="text" class="search-modal-input" placeholder="Search documentation..." autofocus>
                <button class="search-modal-close" aria-label="Close search">
                    <kbd>ESC</kbd>
                </button>
            </div>
            <div class="search-modal-results"></div>
            <div class="search-modal-footer">
                <span><kbd>↑↓</kbd> Navigate</span>
                <span><kbd>↵</kbd> Select</span>
                <span><kbd>ESC</kbd> Close</span>
            </div>
        </div>
    </div>`

var pageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="{{.Page.Title}} - {{.Config.Title}}">
    <title>{{.Page.Title}} | {{.Config.Title}}</title>
    ` + sharedHead + `
    <link rel="stylesheet" href="{{.Config.BaseURL}}assets/styles.css">
</head>
<body>
    <button class="mobile-menu-toggle" aria-label="Toggle navigation menu">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <line x1="3" y1="12" x2="21" y2="12"></line>
            <line x1="3" y1="18" x2="21" y2="18"></line>
        </svg>
    </button>

    <aside class="sidebar">
        <div class="sidebar-header">
            <a href="{{.Config.BaseURL}}index.html" class="logo">
                <span class="logo-icon">` + string(iconMail) + `</span>
                <span class="logo-text">Envelope</span>
            </a>
            ` + themeSwitchHTML + `
        </div>

        ` + searchBlockHTML + `

        <nav class="sidebar-nav">
            {{range .Sections}}
            <div class="nav-section {{if .IsActive}}active{{end}}">
                <button class="nav-section-header" aria-expanded="{{if .IsActive}}true{{else}}false{{end}}">
                    <span class="nav-section-title">{{.Name}}</span>
                    <svg class="nav-section-arrow" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M9 18l6-6-6-6"></path>
                    </svg>
                </button>
                <ul class="nav-section-items">
                    {{$icon := .Icon}}
                    {{range .Pages}}
                    <li>
                        <a href="{{.URL}}" class="nav-link {{if isActive . $.Page}}active{{end}}">
                            <span class="nav-link-icon">{{$icon}}</span>
                            <span>{{.Title}}</span>
                        </a>
                    </li>
                    {{end}}
                </ul>
            </div>
            {{end}}
        </nav>

        ` + sidebarFooterHTML + `
    </aside>

    <main class="main-content">
        <article class="content">
            <div class="content-header">
                <div class="eyebrow">{{.SectionName}}</div>
                <button class="copy-page-button" type="button" data-source="{{.Page.Content}}" aria-label="Copy page as Markdown">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                        <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                    </svg>
                    <span>Copy page</span>
                </button>
            </div>

            <div class="prose">
                {{.Page.HTMLContent}}
            </div>

            <nav class="page-nav">
                <div class="page-nav-info">
                    <span class="page-nav-edit">
                        Found an issue? <a href="https://github.com/envelope-mx/docs/edit/main/docs/{{.Page.Path}}" target="_blank" rel="noopener">Edit this page on GitHub</a>
                    </span>
                </div>
            </nav>
        </article>

        <aside class="toc">
            <div class="toc-header">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="8" y1="6" x2="21" y2="6"></line><line x1="8" y1="12" x2="21" y2="12"></line><line x1="8" y1="18" x2="21" y2="18"></line><line x1="3" y1="6" x2="3.01" y2="6"></line><line x1="3" y1="12" x2="3.01" y2="12"></line><line x1="3" y1="18" x2="3.01" y2="18"></line></svg>
                On this page
            </div>
            <nav class="toc-nav">
                <!-- Populated by JavaScript -->
            </nav>
        </aside>
    </main>

    <div class="sidebar-overlay"></div>

    ` + searchModalHTML + `

    <script src="{{.Config.BaseURL}}assets/app.js"></script>
</body>
</html>`

var indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="{{.Config.Title}} - Self-hosted, API-first, multi-tenant email platform">
    <title>{{.Config.Title}}</title>
    ` + sharedHead + `
    <link rel="stylesheet" href="{{.Config.BaseURL}}assets/styles.css">
</head>
<body>
    <button class="mobile-menu-toggle" aria-label="Toggle navigation menu">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <line x1="3" y1="12" x2="21" y2="12"></line>
            <line x1="3" y1="18" x2="21" y2="18"></line>
        </svg>
    </button>

    <aside class="sidebar">
        <div class="sidebar-header">
            <a href="{{.Config.BaseURL}}index.html" class="logo">
                <span class="logo-icon">` + string(iconMail) + `</span>
                <span class="logo-text">Envelope</span>
            </a>
            ` + themeSwitchHTML + `
        </div>

        ` + searchBlockHTML + `

        <nav class="sidebar-nav">
            {{range .Sections}}
            <div class="nav-section {{if eq .Slug "root"}}active{{end}}">
                <button class="nav-section-header" aria-expanded="{{if eq .Slug "root"}}true{{else}}false{{end}}">
                    <span class="nav-section-title">{{.Name}}</span>
                    <svg class="nav-section-arrow" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M9 18l6-6-6-6"></path>
                    </svg>
                </button>
                <ul class="nav-section-items">
                    {{$icon := .Icon}}
                    {{range .Pages}}
                    <li>
                        <a href="{{.URL}}" class="nav-link">
                            <span class="nav-link-icon">{{$icon}}</span>
                            <span>{{.Title}}</span>
                        </a>
                    </li>
                    {{end}}
                </ul>
            </div>
            {{end}}
        </nav>

        ` + sidebarFooterHTML + `
    </aside>

    <main class="main-content">
        <article class="content home-content">
            <div class="hero">
                <div class="hero-icon">` + string(iconMailLg) + `</div>
                <h1 class="hero-title">Envelope</h1>
                <p class="hero-subtitle">A self-hosted, API-first, multi-tenant email platform. Send and receive mail over a full-featured REST API or authenticated SMTP, with self-serve accounts, automatic DKIM, and signed retryable webhooks — running on your own infrastructure.</p>
                <div class="hero-actions">
                    <a href="{{.Config.BaseURL}}getting-started/quickstart.html" class="btn btn-primary">Quickstart</a>
                    <a href="{{.Config.BaseURL}}api-reference/overview.html" class="btn btn-secondary">API Reference</a>
                </div>
            </div>

            <div class="hero-snippet-wrap">
                <div class="code-block">
                    <div class="code-block-bar">
                        <span class="code-block-dots"><span></span><span></span><span></span></span>
                        <span class="code-block-label">quickstart.sh</span>
                        <button class="copy-button" type="button" aria-label="Copy code">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                            </svg>
                            <span>Copy</span>
                        </button>
                    </div>
                    <pre><code># 1. Create an account (one-time, admin bootstrap token)
curl -X POST https://mail.yourdomain.example/accounts \
  -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"name": "Acme Inc"}'
# → save data.token.token as $ACCOUNT_TOKEN, data.account.id as $ACCOUNT_ID

# 2. Self-serve: create a vhost under that account, no admin token needed from here
curl -X POST https://mail.yourdomain.example/accounts/$ACCOUNT_ID/vhosts \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"domain": "mail.acme.example"}'

# 3. Send your first email over REST
curl -X POST https://mail.yourdomain.example/accounts/$ACCOUNT_ID/messages \
  -H "Authorization: Bearer $ACCOUNT_TOKEN" -H "Content-Type: application/json" \
  -d '{"from": "hello@mail.acme.example", "to": ["you@example.com"],
       "subject": "Hello from Envelope", "text": "It works."}'</code></pre>
                </div>
            </div>

            <div class="features">
                <div class="feature-card">
                    <div class="feature-icon">` + string(iconUsers) + `</div>
                    <h3 class="feature-title">Multi-tenant by default</h3>
                    <p class="feature-description">Accounts own any number of vhosts. Create one account, get a token, and self-serve everything else — vhosts, mailboxes, more tokens — with no further admin involvement.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">` + string(iconSendLg) + `</div>
                    <h3 class="feature-title">REST or SMTP</h3>
                    <p class="feature-description">Send through a full-featured REST endpoint (to/cc/bcc, HTML + text, attachments, custom headers) or authenticated SMTP submission — same DKIM signing, same queue, same webhook.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">` + string(iconLock) + `</div>
                    <h3 class="feature-title">Automatic DKIM</h3>
                    <p class="feature-description">Every vhost gets a 2048-bit RSA key generated the moment it's created. Publish one DNS TXT record and outbound mail is already being signed.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">` + string(iconZapLg) + `</div>
                    <h3 class="feature-title">Signed, retried webhooks</h3>
                    <p class="feature-description">HMAC-signed event delivery for received, queued, delivered, bounced, and deferred mail, with exponential backoff and a queryable dead-letter/redrive history.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">` + string(iconInboxDownLg) + `</div>
                    <h3 class="feature-title">IMAP retrieval</h3>
                    <p class="feature-description">Read inbound and quarantined mail with any standard IMAP4rev2 client — no separate polling integration required.</p>
                </div>
                <div class="feature-card">
                    <div class="feature-icon">` + string(iconAnchorLg) + `</div>
                    <h3 class="feature-title">Deploy anywhere</h3>
                    <p class="feature-description">A single binary or Docker image, running as one process or split across roles behind Docker Compose or Kubernetes.</p>
                </div>
            </div>

            <div class="quick-links">
                <h2 class="quick-links-title">Explore the docs</h2>
                <div class="quick-links-grid">
                    {{range .Sections}}
                    {{if ne .Slug "root"}}
                    <a href="{{if .Pages}}{{(index .Pages 0).URL}}{{else}}#{{end}}" class="quick-link-card">
                        <span class="quick-link-icon">{{.Icon}}</span>
                        <span class="quick-link-title">{{.Name}}</span>
                    </a>
                    {{end}}
                    {{end}}
                </div>
            </div>
        </article>
    </main>

    <div class="sidebar-overlay"></div>

    ` + searchModalHTML + `

    <script src="{{.Config.BaseURL}}assets/app.js"></script>
</body>
</html>`
