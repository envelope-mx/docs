package builder

const cssStyles = `/* ========================================
   Envelope Documentation - Design System
   ======================================== */

/* CSS Custom Properties (Design Tokens) — light palette is the base;
   dark is layered on via prefers-color-scheme (system) and [data-theme]
   (explicit override), per the three-state theme contract. */
:root {
  --color-bg: #ffffff;
  --color-bg-secondary: #f6f5f2;
  --color-bg-tertiary: #ece9e2;
  --color-text: #171512;
  --color-text-secondary: #4a463e;
  --color-text-muted: #837c6d;
  --color-border: #e4e0d6;
  --color-border-light: #ece9e2;
  --color-primary: #9a6b12;
  --color-primary-hover: #7c5510;
  --color-primary-light: #f6ead0;
  --color-accent: #b5461f;
  --color-accent-light: #f8e3d9;
  --color-success: #16733f;
  --color-success-light: #dcf3e4;
  --color-danger: #b3261e;
  --color-danger-light: #fbe2df;
  --color-code-bg: #171512;
  --color-code-text: #ece9e2;

  --sidebar-width: 272px;
  --sidebar-bg: var(--color-bg);
  --toc-width: 220px;

  --font-sans: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  --font-mono: 'JetBrains Mono', 'SF Mono', SFMono-Regular, ui-monospace, 'DejaVu Sans Mono', Menlo, Consolas, monospace;

  --space-xs: 0.25rem;
  --space-sm: 0.5rem;
  --space-md: 1rem;
  --space-lg: 1.5rem;
  --space-xl: 2rem;
  --space-2xl: 3rem;

  --transition-fast: 150ms ease;
  --transition-normal: 250ms ease;

  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.06);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.10);
  --shadow-lg: 0 10px 25px -5px rgba(0, 0, 0, 0.18);

  --radius-sm: 0.375rem;
  --radius-md: 0.5rem;
  --radius-lg: 0.75rem;
  --radius-full: 9999px;
}

/* Dark tokens — applied when the OS prefers dark AND no explicit "light"
   choice has been made ("system", following the OS)... */
@media (prefers-color-scheme: dark) {
  :root:not([data-theme="light"]) {
    --color-bg: #0a0a0a;
    --color-bg-secondary: #131211;
    --color-bg-tertiary: #1d1c19;
    --color-text: #f3f1ec;
    --color-text-secondary: #a8a296;
    --color-text-muted: #6f6a5f;
    --color-border: #262420;
    --color-border-light: #1d1c19;
    --color-primary: #f2b544;
    --color-primary-hover: #f7c976;
    --color-primary-light: #3a2c10;
    --color-accent: #f0703c;
    --color-accent-light: #3a2015;
    --color-success: #34d399;
    --color-success-light: #10281c;
    --color-danger: #f87171;
    --color-danger-light: #331313;
    --color-code-bg: #050505;
    --color-code-text: #ece9e2;
    --sidebar-bg: #0a0a0a;
    --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.4);
    --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.5);
    --shadow-lg: 0 10px 25px -5px rgba(0, 0, 0, 0.6);
  }
}

/* ...and again when the visitor has explicitly chosen dark, so the
   toggle always wins over the OS setting in both directions. */
:root[data-theme="dark"] {
  --color-bg: #0a0a0a;
  --color-bg-secondary: #131211;
  --color-bg-tertiary: #1d1c19;
  --color-text: #f3f1ec;
  --color-text-secondary: #a8a296;
  --color-text-muted: #6f6a5f;
  --color-border: #262420;
  --color-border-light: #1d1c19;
  --color-primary: #f2b544;
  --color-primary-hover: #f7c976;
  --color-primary-light: #3a2c10;
  --color-accent: #f0703c;
  --color-accent-light: #3a2015;
  --color-success: #34d399;
  --color-success-light: #10281c;
  --color-danger: #f87171;
  --color-danger-light: #331313;
  --color-code-bg: #050505;
  --color-code-text: #ece9e2;
  --sidebar-bg: #0a0a0a;
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.4);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.5);
  --shadow-lg: 0 10px 25px -5px rgba(0, 0, 0, 0.6);
}

/* Reset & Base */
*, *::before, *::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

html {
  font-size: 16px;
  scroll-behavior: smooth;
  -webkit-text-size-adjust: 100%;
}

body {
  font-family: var(--font-sans);
  font-size: 1rem;
  line-height: 1.65;
  color: var(--color-text);
  background-color: var(--color-bg);
  display: flex;
  min-height: 100vh;
  overflow-x: hidden;
}

::selection {
  background: var(--color-primary-light);
  color: var(--color-primary);
}

/* ========================================
   Sidebar
   ======================================== */

.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: var(--sidebar-width);
  background: var(--sidebar-bg);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  z-index: 100;
  transition: transform var(--transition-normal);
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-lg) var(--space-lg) var(--space-md);
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  text-decoration: none;
  color: var(--color-text);
  font-weight: 700;
  font-size: 1.1rem;
  letter-spacing: -0.01em;
}

.logo-icon {
  display: flex;
  color: var(--color-primary);
}

.logo:hover {
  color: var(--color-primary);
}

/* Theme switch — three-way segmented control */
.theme-switch {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 2px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-full);
  background: var(--color-bg-secondary);
}

.theme-switch-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border: none;
  border-radius: var(--radius-full);
  background: none;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.theme-switch-btn:hover {
  color: var(--color-text);
}

.theme-switch-btn.active {
  background: var(--color-bg-tertiary);
  color: var(--color-primary);
}

/* Search */
.search-container {
  padding: 0 var(--space-lg) var(--space-md);
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: var(--space-md);
  color: var(--color-text-muted);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 0.55rem var(--space-md);
  padding-left: 2.35rem;
  padding-right: 3rem;
  font-family: var(--font-sans);
  font-size: 0.875rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-full);
  background: var(--color-bg-secondary);
  color: var(--color-text);
  transition: all var(--transition-fast);
}

.search-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px var(--color-primary-light);
}

.search-input::placeholder {
  color: var(--color-text-muted);
}

.search-shortcut {
  position: absolute;
  right: var(--space-sm);
  padding: 2px 6px;
  font-size: 0.75rem;
  font-family: var(--font-sans);
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-muted);
  pointer-events: none;
}

/* Sidebar Nav */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-sm) var(--space-md) var(--space-lg);
}

.sidebar-nav::-webkit-scrollbar {
  width: 6px;
}

.sidebar-nav::-webkit-scrollbar-thumb {
  background: var(--color-border);
  border-radius: var(--radius-full);
}

.nav-section {
  margin-bottom: 2px;
}

.nav-section-header {
  display: flex;
  align-items: center;
  width: 100%;
  padding: var(--space-sm) var(--space-sm);
  background: none;
  border: none;
  border-radius: var(--radius-md);
  color: var(--color-text);
  font-family: var(--font-sans);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast);
  text-align: left;
}

.nav-section-header:hover {
  background: var(--color-bg-secondary);
}

.nav-section-title {
  flex: 1;
}

.nav-section-arrow {
  transition: transform var(--transition-fast);
  transform: rotate(90deg);
  color: var(--color-text-muted);
}

.nav-section:not(.active) .nav-section-arrow {
  transform: rotate(0deg);
}

.nav-section-items {
  list-style: none;
  overflow: hidden;
  max-height: 0;
  transition: max-height var(--transition-normal);
}

.nav-section.active .nav-section-items {
  max-height: 1200px;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.4rem var(--space-sm) 0.4rem calc(var(--space-sm) + 0.5rem);
  margin: 1px 0;
  border-radius: 0 var(--radius-md) var(--radius-md) 0;
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 0.875rem;
  border-left: 2px solid transparent;
  transition: all var(--transition-fast);
}

.nav-link-icon {
  display: flex;
  flex-shrink: 0;
  color: var(--color-text-muted);
  transition: color var(--transition-fast);
}

.nav-link:hover {
  color: var(--color-text);
}

.nav-link:hover .nav-link-icon {
  color: var(--color-text);
}

.nav-link.active {
  color: var(--color-primary);
  font-weight: 600;
  border-left-color: var(--color-primary);
  background: var(--color-bg-secondary);
}

.nav-link.active .nav-link-icon {
  color: var(--color-primary);
}

.sidebar-footer {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: var(--space-md) var(--space-lg);
  border-top: 1px solid var(--color-border);
}

.sidebar-footer-link {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
  padding: 0.3rem 0;
  color: var(--color-text-muted);
  text-decoration: none;
  font-size: 0.8125rem;
  transition: color var(--transition-fast);
}

.sidebar-footer-link:hover {
  color: var(--color-text);
}

/* ========================================
   Main Content Layout
   ======================================== */

.main-content {
  margin-left: var(--sidebar-width);
  flex: 1;
  display: flex;
  justify-content: center;
  min-width: 0;
}

.content {
  width: 100%;
  max-width: 860px;
  padding: var(--space-2xl) var(--space-xl) 6rem;
}

.content-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-md);
  margin-bottom: 0.5rem;
}

.eyebrow {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-primary);
}

.copy-page-button {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.4rem 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-full);
  background: var(--color-bg-secondary);
  color: var(--color-text-secondary);
  font-family: var(--font-sans);
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
  flex-shrink: 0;
  margin-top: 0.25rem;
}

.copy-page-button:hover {
  border-color: var(--color-primary);
  color: var(--color-text);
}

.copy-page-button.copied {
  color: var(--color-success);
  border-color: var(--color-success);
}

/* ========================================
   Prose (Markdown Content)
   ======================================== */

.prose {
  color: var(--color-text);
}

.prose h1, .prose h2, .prose h3, .prose h4, .prose h5, .prose h6 {
  font-weight: 800;
  letter-spacing: -0.015em;
  margin-top: 2.5rem;
  margin-bottom: 1rem;
  scroll-margin-top: 2rem;
}

.prose h1 { font-size: 2.25rem; margin-top: 0.25rem; }
.prose h2 { font-size: 1.625rem; border-bottom: 1px solid var(--color-border); padding-bottom: var(--space-sm); }
.prose h3 { font-size: 1.25rem; }
.prose h4 { font-size: 1.0625rem; }

.prose p {
  margin-bottom: 1.15rem;
  color: var(--color-text-secondary);
}

.prose a {
  color: var(--color-text);
  text-decoration: underline;
  text-decoration-color: var(--color-border);
  text-underline-offset: 2px;
}

.prose a:hover {
  color: var(--color-primary);
  text-decoration-color: var(--color-primary);
}

.prose strong {
  color: var(--color-text);
  font-weight: 700;
}

.prose ul, .prose ol {
  margin: 0 0 1.15rem 1.5rem;
  color: var(--color-text-secondary);
}

.prose li {
  margin-bottom: 0.4rem;
}

.prose li::marker {
  color: var(--color-primary);
}

.prose code {
  font-family: var(--font-mono);
  font-size: 0.85em;
  background: var(--color-bg-tertiary);
  color: var(--color-text);
  padding: 0.15em 0.4em;
  border-radius: var(--radius-sm);
}

.prose pre {
  font-family: var(--font-mono);
  font-size: 0.85rem;
  line-height: 1.6;
  background: var(--color-code-bg);
  color: var(--color-code-text);
  padding: var(--space-lg);
  border-radius: var(--radius-lg);
  overflow-x: auto;
  margin-bottom: 1.15rem;
}

.prose pre code {
  background: none;
  padding: 0;
  color: inherit;
  font-size: 1em;
}

.prose table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 1.15rem;
  font-size: 0.9rem;
  display: block;
  overflow-x: auto;
}

.prose th, .prose td {
  padding: 0.6rem 0.85rem;
  border: 1px solid var(--color-border);
  text-align: left;
}

.prose th {
  background: var(--color-bg-secondary);
  font-weight: 700;
  color: var(--color-text);
  white-space: nowrap;
}

.prose tr:nth-child(even) {
  background: var(--color-bg-secondary);
}

.prose blockquote {
  border-left: 3px solid var(--color-primary);
  padding: 0.25rem 0 0.25rem 1.15rem;
  margin: 0 0 1.15rem;
  color: var(--color-text-muted);
}

.prose blockquote p:last-child {
  margin-bottom: 0;
}

.prose hr {
  border: none;
  border-top: 1px solid var(--color-border);
  margin: 2.5rem 0;
}

.prose img {
  max-width: 100%;
  border-radius: var(--radius-md);
}

.prose input[type="checkbox"] {
  margin-right: 0.4em;
}

/* ========================================
   Code block chrome (bar + copy button)
   ======================================== */

.code-block {
  position: relative;
  margin-bottom: 1.15rem;
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--color-border);
  background: var(--color-code-bg);
}

.code-block .prose,
.code-block pre {
  margin-bottom: 0;
  border-radius: 0;
}

.code-block-bar {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: 0.55rem var(--space-md);
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.code-block-dots {
  display: flex;
  gap: 5px;
}

.code-block-dots span {
  width: 9px;
  height: 9px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.14);
}

.code-block-label {
  flex: 1;
  font-family: var(--font-mono);
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.4);
}

.copy-button {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.3rem 0.6rem;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-sm);
  color: rgba(255, 255, 255, 0.7);
  font-family: var(--font-sans);
  font-size: 0.75rem;
  cursor: pointer;
  transition: all var(--transition-fast);
  margin-left: auto;
}

.copy-button:hover {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.copy-button.copied {
  color: var(--color-success);
  border-color: var(--color-success);
}

/* ========================================
   Table of Contents
   ======================================== */

.toc {
  display: none;
  width: var(--toc-width);
  flex-shrink: 0;
  padding: var(--space-2xl) var(--space-lg) 0 0;
}

@media (min-width: 1280px) {
  .toc { display: block; }
}

.toc-header {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.8125rem;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--space-sm);
}

.toc-nav {
  position: sticky;
  top: var(--space-2xl);
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  border-left: 1px solid var(--color-border);
  padding-left: var(--space-md);
  max-height: calc(100vh - 4rem);
  overflow-y: auto;
}

.toc-nav a {
  position: relative;
  color: var(--color-text-muted);
  text-decoration: none;
  font-size: 0.85rem;
  line-height: 1.4;
}

.toc-nav a::before {
  content: "";
  position: absolute;
  left: -1.0625rem;
  top: 0;
  bottom: 0;
  width: 2px;
  background: transparent;
}

.toc-nav a:hover {
  color: var(--color-text);
}

.toc-nav a.active {
  color: var(--color-primary);
  font-weight: 600;
}

.toc-nav a.active::before {
  background: var(--color-primary);
}

.toc-nav a.toc-h3 {
  padding-left: var(--space-md);
  font-size: 0.8rem;
}

/* ========================================
   Page nav / edit link
   ======================================== */

.page-nav {
  margin-top: var(--space-2xl);
  padding-top: var(--space-lg);
  border-top: 1px solid var(--color-border);
}

.page-nav-info {
  font-size: 0.875rem;
  color: var(--color-text-muted);
}

.page-nav-edit a {
  color: var(--color-primary);
  text-decoration: none;
}

.page-nav-edit a:hover {
  text-decoration: underline;
}

/* ========================================
   Callouts / admonitions
   ======================================== */

.callout {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: var(--space-md) var(--space-lg);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  background: var(--color-bg-secondary);
  margin-bottom: 1.15rem;
  font-size: 0.9rem;
  color: var(--color-text-secondary);
}

.callout::before {
  display: flex;
  width: 22px;
  height: 22px;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  font-size: 0.9rem;
}

.callout p:last-child {
  margin-bottom: 0;
}

.callout.note::before {
  content: "i";
  font-weight: 700;
  font-style: italic;
  font-family: Georgia, serif;
  background: var(--color-primary-light);
  color: var(--color-primary);
}

.callout.warning::before {
  content: "!";
  font-weight: 800;
  background: var(--color-accent-light);
  color: var(--color-accent);
}

.callout.gap {
  border-style: dashed;
}

.callout.gap::before {
  content: "◔";
  background: var(--color-bg-tertiary);
  color: var(--color-text-muted);
}

/* ========================================
   HTTP method badges
   ======================================== */

.method-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.15em 0.55em;
  border-radius: var(--radius-sm);
  font-family: var(--font-mono);
  font-size: 0.75em;
  font-weight: 700;
  letter-spacing: 0.02em;
  color: #fff;
  vertical-align: middle;
}

.method-badge.get { background: #2563eb; }
.method-badge.post { background: #16a34a; }
.method-badge.patch { background: #d97706; }
.method-badge.delete { background: #dc2626; }

/* ========================================
   Home Page
   ======================================== */

.home-content {
  max-width: 1000px;
}

.hero {
  text-align: center;
  padding: var(--space-2xl) 0 var(--space-xl);
}

.hero-icon {
  display: flex;
  justify-content: center;
  color: var(--color-primary);
  margin-bottom: var(--space-md);
}

.hero-title {
  font-size: 3.25rem;
  font-weight: 800;
  letter-spacing: -0.03em;
  margin-bottom: var(--space-md);
  color: var(--color-text);
}

.hero-subtitle {
  font-size: 1.15rem;
  line-height: 1.6;
  color: var(--color-text-secondary);
  max-width: 640px;
  margin: 0 auto var(--space-xl);
}

.hero-actions {
  display: flex;
  gap: var(--space-md);
  justify-content: center;
}

.btn {
  display: inline-flex;
  align-items: center;
  padding: 0.7rem 1.4rem;
  border-radius: var(--radius-full);
  font-weight: 600;
  font-size: 0.95rem;
  text-decoration: none;
  transition: all var(--transition-fast);
}

.btn-primary {
  background: var(--color-primary);
  color: #171512;
}

.btn-primary:hover {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.btn-secondary {
  background: var(--color-bg-secondary);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.btn-secondary:hover {
  background: var(--color-bg-tertiary);
}

.hero-snippet-wrap {
  max-width: 760px;
  margin: 0 auto var(--space-2xl);
}

.hero-snippet-wrap .code-block-label {
  color: rgba(255, 255, 255, 0.45);
}

.features {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-lg);
  margin-bottom: var(--space-2xl);
}

.feature-card {
  padding: var(--space-lg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-bg-secondary);
  transition: all var(--transition-fast);
}

.feature-card:hover {
  border-color: var(--color-primary);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.feature-icon {
  display: flex;
  color: var(--color-primary);
  margin-bottom: var(--space-md);
}

.feature-title {
  font-size: 1.05rem;
  font-weight: 700;
  margin-bottom: 0.4rem;
}

.feature-description {
  font-size: 0.9rem;
  color: var(--color-text-muted);
  line-height: 1.55;
}

.quick-links {
  margin-bottom: var(--space-xl);
}

.quick-links-title {
  font-size: 1.1rem;
  font-weight: 700;
  margin-bottom: var(--space-md);
  text-align: center;
}

.quick-links-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-sm);
}

.quick-link-card {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-md);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  text-decoration: none;
  color: var(--color-text);
  font-weight: 600;
  font-size: 0.9rem;
  transition: all var(--transition-fast);
}

.quick-link-card:hover {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
  color: var(--color-primary);
}

.quick-link-icon {
  display: flex;
  color: var(--color-text-muted);
}

.quick-link-card:hover .quick-link-icon {
  color: var(--color-primary);
}

/* ========================================
   Search Modal
   ======================================== */

.search-modal {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 1000;
  align-items: flex-start;
  justify-content: center;
  padding-top: 12vh;
}

.search-modal.active {
  display: flex;
}

.search-modal-content {
  width: 100%;
  max-width: 600px;
  max-height: 70vh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
  margin: 0 var(--space-md);
}

.search-modal-header {
  display: flex;
  align-items: center;
  gap: var(--space-sm);
  padding: var(--space-md) var(--space-lg);
  border-bottom: 1px solid var(--color-border);
}

.search-modal-header .search-icon {
  position: static;
  color: var(--color-text-muted);
}

.search-modal-input {
  flex: 1;
  border: none;
  outline: none;
  background: none;
  font-family: var(--font-sans);
  font-size: 1rem;
  color: var(--color-text);
}

.search-modal-input::placeholder {
  color: var(--color-text-muted);
}

.search-modal-close {
  border: none;
  background: none;
  cursor: pointer;
}

.search-modal-close kbd {
  padding: 2px 6px;
  font-size: 0.75rem;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  color: var(--color-text-muted);
}

.search-modal-results {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-sm);
}

.search-result-item {
  display: block;
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--radius-md);
  text-decoration: none;
  color: var(--color-text);
  margin-bottom: 2px;
}

.search-result-item.selected,
.search-result-item:hover {
  background: var(--color-bg-secondary);
}

.search-result-section {
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  color: var(--color-primary);
  margin-bottom: 0.15rem;
}

.search-result-title {
  font-weight: 600;
  margin-bottom: 0.15rem;
}

.search-result-excerpt {
  font-size: 0.8rem;
  color: var(--color-text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.search-result-item mark {
  background: var(--color-accent-light);
  color: var(--color-accent);
  border-radius: 2px;
  padding: 0 1px;
}

.search-no-results {
  padding: var(--space-lg);
  text-align: center;
  color: var(--color-text-muted);
  font-size: 0.9rem;
}

.search-modal-footer {
  display: flex;
  gap: var(--space-md);
  padding: var(--space-sm) var(--space-lg);
  border-top: 1px solid var(--color-border);
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.search-modal-footer kbd {
  padding: 1px 5px;
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  margin-right: 0.25rem;
}

/* ========================================
   Mobile menu / overlay
   ======================================== */

.mobile-menu-toggle {
  display: none;
  position: fixed;
  top: var(--space-md);
  left: var(--space-md);
  z-index: 200;
  width: 40px;
  height: 40px;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text);
  cursor: pointer;
  box-shadow: var(--shadow-md);
}

.sidebar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 99;
}

.sidebar-overlay.active {
  display: block;
}

/* ========================================
   Responsive
   ======================================== */

@media (max-width: 1024px) {
  .mobile-menu-toggle {
    display: flex;
  }

  .sidebar {
    transform: translateX(-100%);
    box-shadow: var(--shadow-lg);
  }

  .sidebar.open {
    transform: translateX(0);
  }

  .main-content {
    margin-left: 0;
    padding-top: 3.5rem;
  }

  .toc {
    display: none !important;
  }
}

@media (max-width: 900px) {
  .features {
    grid-template-columns: repeat(2, 1fr);
  }

  .quick-links-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .content {
    padding: var(--space-xl) var(--space-md) 4rem;
  }

  .content-header {
    flex-direction: column;
  }

  .hero-title {
    font-size: 2.25rem;
  }

  .hero-subtitle {
    font-size: 1rem;
  }

  .hero-actions {
    flex-direction: column;
  }

  .features {
    grid-template-columns: 1fr;
  }

  .quick-links-grid {
    grid-template-columns: 1fr;
  }

  .prose h1 { font-size: 1.75rem; }
  .prose h2 { font-size: 1.4rem; }
}

/* ========================================
   Chroma syntax highlighting (Envelope theme)
   ======================================== */
.chroma { color: #ece9e2; }
.chroma .err { color: #f87171; }
.chroma .lntd { vertical-align: top; padding: 0; margin: 0; border: 0; }
.chroma .lntable { border-spacing: 0; padding: 0; margin: 0; border: 0; width: auto; overflow: auto; display: block; }
.chroma .hl { display: block; width: 100%; background-color: #1d1c19; }
.chroma .lnt { margin-right: 0.4em; padding: 0 0.4em; color: #6f6a5f; }
.chroma .ln { margin-right: 0.4em; padding: 0 0.4em; color: #6f6a5f; }
.chroma .k { color: #f2b544; }
.chroma .kc { color: #f2b544; }
.chroma .kd { color: #6bc9ff; font-style: italic; }
.chroma .kn { color: #f2b544; }
.chroma .kp { color: #f2b544; }
.chroma .kr { color: #f2b544; }
.chroma .kt { color: #6bc9ff; }
.chroma .na { color: #7fd9a8; }
.chroma .nb { color: #6bc9ff; font-style: italic; }
.chroma .nc { color: #7fd9a8; }
.chroma .no { color: #f0703c; }
.chroma .nd { color: #7fd9a8; }
.chroma .ne { color: #7fd9a8; }
.chroma .nf { color: #7fd9a8; }
.chroma .nl { color: #6bc9ff; font-style: italic; }
.chroma .nn { color: #ece9e2; }
.chroma .nt { color: #f2b544; }
.chroma .nv { color: #6bc9ff; font-style: italic; }
.chroma .s { color: #a8d977; }
.chroma .sa { color: #a8d977; }
.chroma .sb { color: #a8d977; }
.chroma .sc { color: #a8d977; }
.chroma .dl { color: #a8d977; }
.chroma .sd { color: #a8d977; }
.chroma .s2 { color: #a8d977; }
.chroma .se { color: #f0703c; }
.chroma .sh { color: #a8d977; }
.chroma .si { color: #f0703c; }
.chroma .sx { color: #a8d977; }
.chroma .sr { color: #7fd9a8; }
.chroma .s1 { color: #a8d977; }
.chroma .ss { color: #f2b544; }
.chroma .m { color: #f0703c; }
.chroma .mb { color: #f0703c; }
.chroma .mf { color: #f0703c; }
.chroma .mh { color: #f0703c; }
.chroma .mi { color: #f0703c; }
.chroma .il { color: #f0703c; }
.chroma .mo { color: #f0703c; }
.chroma .o { color: #f2b544; }
.chroma .ow { color: #f2b544; }
.chroma .p { color: #a8a296; }
.chroma .c { color: #6f6a5f; font-style: italic; }
.chroma .ch { color: #6f6a5f; font-style: italic; }
.chroma .cm { color: #6f6a5f; font-style: italic; }
.chroma .cp { color: #6f6a5f; }
.chroma .cpf { color: #6f6a5f; }
.chroma .c1 { color: #6f6a5f; font-style: italic; }
.chroma .cs { color: #6f6a5f; font-style: italic; }
.chroma .gd { color: #f87171; }
.chroma .ge { font-style: italic; }
.chroma .gr { color: #f87171; }
.chroma .gh { color: #6bc9ff; font-weight: bold; }
.chroma .gi { color: #7fd9a8; }
.chroma .go { color: #a8a296; }
.chroma .gp { color: #6f6a5f; }
.chroma .gs { font-weight: bold; }
.chroma .gu { color: #f2b544; font-weight: bold; }
.chroma .gt { color: #f87171; }
.chroma .w { color: #6f6a5f; }
.chroma .bp { color: #6bc9ff; }
.chroma .fm { color: #7fd9a8; }
.chroma .vc { color: #6bc9ff; }
.chroma .vg { color: #6bc9ff; }
.chroma .vi { color: #6bc9ff; }
.chroma .vm { color: #6bc9ff; }
`
