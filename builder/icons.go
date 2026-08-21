package builder

import "html/template"

// Line icons used throughout the generated site — a consistent 24x24
// stroke-based icon set (matching the sidebar's search/menu/theme icons),
// used in place of emoji for section navigation, feature highlights, and
// UI chrome. Each is deliberately built from simple primitives (circles,
// rects, lines, short paths) rather than reproduced from a specific
// third-party icon library's exact path data.

const iconAttrs = `fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"`

func svgIcon(size, body string) template.HTML {
	return template.HTML(`<svg width="` + size + `" height="` + size + `" viewBox="0 0 24 24" ` + iconAttrs + `>` + body + `</svg>`)
}

var (
	iconMail = svgIcon("16", `<rect x="3" y="5" width="18" height="14" rx="2"></rect><polyline points="3 7 12 13 21 7"></polyline>`)

	iconArrowUpRight = svgIcon("16", `<line x1="7" y1="17" x2="17" y2="7"></line><polyline points="7 7 17 7 17 17"></polyline>`)

	iconAnchor = svgIcon("16", `<circle cx="12" cy="5" r="3"></circle><line x1="12" y1="22" x2="12" y2="8"></line><path d="M5 12H2a10 10 0 0 0 20 0h-3"></path>`)

	iconCompass = svgIcon("16", `<circle cx="12" cy="12" r="10"></circle><polygon points="16.24 7.76 14.12 14.12 7.76 16.24 9.88 9.88 16.24 7.76"></polygon>`)

	iconCode = svgIcon("16", `<polyline points="16 18 22 12 16 6"></polyline><polyline points="8 6 2 12 8 18"></polyline>`)

	iconZap = svgIcon("16", `<polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>`)

	iconSend = svgIcon("16", `<line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>`)

	iconInboxDown = svgIcon("16", `<path d="M4 4h16v6H4z"></path><polyline points="4 14 12 20 20 14"></polyline><line x1="12" y1="10" x2="12" y2="20"></line>`)

	iconBookOpen = svgIcon("16", `<path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>`)

	iconBookmark = svgIcon("16", `<path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z"></path>`)

	iconFile = svgIcon("16", `<path d="M14 3v5a2 2 0 0 0 2 2h5"></path><path d="M6 3h8l6 6v10a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2z"></path>`)

	iconUsers = svgIcon("22", `<circle cx="8" cy="10" r="4"></circle><circle cx="16" cy="10" r="4"></circle>`)

	iconLock = svgIcon("22", `<rect x="3" y="11" width="18" height="11" rx="2"></rect><path d="M7 11V7a5 5 0 0 1 10 0v4"></path>`)

	iconAnchorLg = svgIcon("22", `<circle cx="12" cy="5" r="3"></circle><line x1="12" y1="22" x2="12" y2="8"></line><path d="M5 12H2a10 10 0 0 0 20 0h-3"></path>`)

	iconZapLg = svgIcon("22", `<polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>`)

	iconInboxDownLg = svgIcon("22", `<path d="M4 4h16v6H4z"></path><polyline points="4 14 12 20 20 14"></polyline><line x1="12" y1="10" x2="12" y2="20"></line>`)

	iconSendLg = svgIcon("22", `<line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>`)

	iconMailLg = svgIcon("40", `<rect x="3" y="5" width="18" height="14" rx="2"></rect><polyline points="3 7 12 13 21 7"></polyline>`)
)

// sectionIcons maps each top-level docs section to its sidebar/nav icon.
var sectionIcons = map[string]template.HTML{
	"root":            iconMail,
	"getting-started": iconArrowUpRight,
	"deployment":      iconAnchor,
	"core-concepts":   iconCompass,
	"api-reference":   iconCode,
	"webhooks":        iconZap,
	"smtp":            iconSend,
	"imap":            iconInboxDown,
	"guides":          iconBookOpen,
	"reference":       iconBookmark,
}
