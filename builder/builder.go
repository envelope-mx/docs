package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// mdLinkTransformer rewrites internal cross-reference links written as
// "../section/page.md" (the natural way to link between source files, and
// how GitHub itself renders them) into "../section/page.html" — the actual
// filename this builder writes to dist/. Without it, every prose link
// between docs pages 404s in the built site while still looking correct in
// the raw Markdown source.
type mdLinkTransformer struct{}

func (t *mdLinkTransformer) Transform(doc *ast.Document, _ text.Reader, _ parser.Context) {
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}
		dest := string(link.Destination)
		if dest == "" || strings.Contains(dest, "://") || strings.HasPrefix(dest, "mailto:") || strings.HasPrefix(dest, "#") {
			return ast.WalkContinue, nil
		}
		path, frag, hasFrag := strings.Cut(dest, "#")
		if strings.HasSuffix(path, ".md") {
			path = strings.TrimSuffix(path, ".md") + ".html"
		}
		if hasFrag {
			path = path + "#" + frag
		}
		link.Destination = []byte(path)
		return ast.WalkContinue, nil
	})
}

// Config holds the build configuration
type Config struct {
	InputDir  string
	OutputDir string
	Title     string
	BaseURL   string
}

// Page represents a documentation page
type Page struct {
	Title       string `json:"title"`
	Path        string `json:"path"`
	URL         string `json:"url"`
	Content     string `json:"content,omitempty"`
	HTMLContent template.HTML
	Section     string `json:"section"`
	Order       int    `json:"order"`
}

// Section represents a documentation section
type Section struct {
	Name  string
	Slug  string
	Pages []Page
	Icon  template.HTML
}

// NavItem represents a navigation item
type NavItem struct {
	Name     string
	Slug     string
	Pages    []Page
	Icon     template.HTML
	IsActive bool
}

// SearchIndex represents the search index
type SearchIndex struct {
	Pages []SearchPage `json:"pages"`
}

// SearchPage represents a page in the search index
type SearchPage struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Section string `json:"section"`
	Content string `json:"content"`
}

// Build builds the documentation site
func Build(config Config) error {
	// Normalize config values with defaults
	if config.InputDir == "" {
		config.InputDir = "docs"
	}
	if config.OutputDir == "" {
		config.OutputDir = "dist"
	}
	if config.Title == "" {
		config.Title = "Envelope Documentation"
	}
	if config.BaseURL == "" {
		config.BaseURL = "/"
	}
	// Ensure BaseURL ends with /
	if !strings.HasSuffix(config.BaseURL, "/") {
		config.BaseURL = config.BaseURL + "/"
	}

	// Initialize goldmark with syntax highlighting
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(util.Prioritized(&mdLinkTransformer{}, 100)),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	// Collect all markdown files
	pages, err := collectPages(config.InputDir, config.BaseURL, md)
	if err != nil {
		return fmt.Errorf("failed to collect pages: %w", err)
	}

	// Organize pages into sections
	sections := organizeSections(pages)

	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate HTML pages
	for i := range pages {
		if err := generatePage(&pages[i], sections, config); err != nil {
			return fmt.Errorf("failed to generate page %s: %w", pages[i].Path, err)
		}
	}

	// Generate index page
	if err := generateIndexPage(sections, config); err != nil {
		return fmt.Errorf("failed to generate index page: %w", err)
	}

	// Generate search index
	if err := generateSearchIndex(pages, config); err != nil {
		return fmt.Errorf("failed to generate search index: %w", err)
	}

	// Copy static assets
	if err := copyAssets(config.OutputDir); err != nil {
		return fmt.Errorf("failed to copy assets: %w", err)
	}

	return nil
}

func collectPages(inputDir, baseURL string, md goldmark.Markdown) ([]Page, error) {
	var pages []Page

	err := filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Convert markdown to HTML
		var buf bytes.Buffer
		if err := md.Convert(content, &buf); err != nil {
			return fmt.Errorf("failed to convert markdown %s: %w", path, err)
		}

		// Extract title from content
		title := extractTitle(string(content))
		if title == "" {
			title = strings.TrimSuffix(filepath.Base(path), ".md")
			title = toTitle(strings.ReplaceAll(title, "-", " "))
		}

		// Calculate relative path
		relPath, _ := filepath.Rel(inputDir, path)
		section := filepath.Dir(relPath)
		if section == "." {
			section = "root"
		}

		// Generate URL — always root-relative (leading "/"), never a bare
		// relative path. A bare "getting-started/introduction.html" would
		// resolve incorrectly from a page nested in a different section
		// (e.g. "/imap/overview.html" would resolve it against "/imap/",
		// not the site root) — every nav link is rendered on every page
		// regardless of that page's own directory depth, so it must be
		// unambiguous from any location.
		url := strings.TrimSuffix(relPath, ".md") + ".html"
		url = filepath.ToSlash(url)
		url = strings.TrimSuffix(baseURL, "/") + "/" + url

		// Determine order based on filename or position
		order := getPageOrder(section, filepath.Base(path))

		pages = append(pages, Page{
			Title:       title,
			Path:        relPath,
			URL:         url,
			Content:     string(content),
			HTMLContent: template.HTML(buf.String()),
			Section:     section,
			Order:       order,
		})

		return nil
	})

	return pages, err
}

func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}

func toTitle(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// pageOrderMaps gives each section its own filename -> order table, since
// the same filename (e.g. "overview.md") is reused as the first page of
// several unrelated sections with completely different meanings.
var pageOrderMaps = map[string]map[string]int{
	"root": {
		"index.md": 0,
	},
	"getting-started": {
		"introduction.md":   1,
		"quickstart.md":     2,
		"authentication.md": 3,
	},
	"deployment": {
		"deploy.md": 1,
	},
	"core-concepts": {
		"multi-tenancy.md":          1,
		"dns-and-dkim.md":           2,
		"rate-limits-and-quotas.md": 3,
		"errors-and-responses.md":   4,
	},
	"api-reference": {
		"overview.md":    1,
		"accounts.md":    2,
		"vhosts.md":      3,
		"mailboxes.md":   4,
		"tokens.md":      5,
		"messages.md":    6,
		"webhooks.md":    7,
		"audit-log.md":   8,
		"data-export.md": 9,
	},
	"webhooks": {
		"overview.md":   1,
		"events.md":     2,
		"signatures.md": 3,
	},
	"smtp": {
		"inbound.md":    1,
		"submission.md": 2,
	},
	"imap": {
		"overview.md":    1,
		"limitations.md": 2,
	},
	"guides": {
		"send-via-rest.md":                1,
		"send-via-smtp.md":                2,
		"receive-and-verify-webhooks.md":  3,
		"handle-bounces-and-deferrals.md": 4,
		"rotate-a-token.md":               5,
		"set-up-a-new-vhost.md":           6,
	},
	"reference": {
		"status-codes.md":      1,
		"glossary.md":          2,
		"known-limitations.md": 3,
	},
}

func getPageOrder(section, filename string) int {
	if m, ok := pageOrderMaps[section]; ok {
		if order, ok := m[filename]; ok {
			return order
		}
	}
	return 100
}

func organizeSections(pages []Page) []Section {
	sectionMap := make(map[string]*Section)
	sectionOrder := map[string]int{
		"root":            0,
		"getting-started": 1,
		"deployment":      2,
		"core-concepts":   3,
		"api-reference":   4,
		"webhooks":        5,
		"smtp":            6,
		"imap":            7,
		"guides":          8,
		"reference":       9,
	}

	sectionNames := map[string]string{
		"root":            "Overview",
		"getting-started": "Getting Started",
		"deployment":      "Deployment",
		"core-concepts":   "Core Concepts",
		"api-reference":   "API Reference",
		"webhooks":        "Webhooks",
		"smtp":            "SMTP",
		"imap":            "IMAP",
		"guides":          "Guides",
		"reference":       "Reference",
	}

	for _, page := range pages {
		if _, ok := sectionMap[page.Section]; !ok {
			name := sectionNames[page.Section]
			if name == "" {
				name = toTitle(strings.ReplaceAll(page.Section, "-", " "))
			}
			icon := sectionIcons[page.Section]
			if icon == "" {
				icon = iconFile
			}
			sectionMap[page.Section] = &Section{
				Name:  name,
				Slug:  page.Section,
				Icon:  icon,
				Pages: []Page{},
			}
		}
		sectionMap[page.Section].Pages = append(sectionMap[page.Section].Pages, page)
	}

	var sections []Section
	for _, section := range sectionMap {
		sort.Slice(section.Pages, func(i, j int) bool {
			if section.Pages[i].Order != section.Pages[j].Order {
				return section.Pages[i].Order < section.Pages[j].Order
			}
			return section.Pages[i].Title < section.Pages[j].Title
		})
		sections = append(sections, *section)
	}

	sort.Slice(sections, func(i, j int) bool {
		oi := sectionOrder[sections[i].Slug]
		oj := sectionOrder[sections[j].Slug]
		if oi != oj {
			return oi < oj
		}
		return sections[i].Name < sections[j].Name
	})

	return sections
}

func generatePage(page *Page, sections []Section, config Config) error {
	var navItems []NavItem
	var sectionName string
	for _, section := range sections {
		isActive := section.Slug == page.Section
		if isActive {
			sectionName = section.Name
		}
		navItems = append(navItems, NavItem{
			Name:     section.Name,
			Slug:     section.Slug,
			Pages:    section.Pages,
			Icon:     section.Icon,
			IsActive: isActive,
		})
	}

	data := struct {
		Page        *Page
		Sections    []NavItem
		Config      Config
		SectionName string
	}{
		Page:        page,
		Sections:    navItems,
		Config:      config,
		SectionName: sectionName,
	}

	tmpl, err := template.New("page").Funcs(template.FuncMap{
		"isActive": func(p Page, current *Page) bool {
			return p.Path == current.Path
		},
	}).Parse(pageTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	outputPath := filepath.Join(config.OutputDir, strings.TrimSuffix(page.Path, ".md")+".html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func generateIndexPage(sections []Section, config Config) error {
	var navItems []NavItem
	for _, section := range sections {
		navItems = append(navItems, NavItem{
			Name:     section.Name,
			Slug:     section.Slug,
			Pages:    section.Pages,
			Icon:     section.Icon,
			IsActive: section.Slug == "root",
		})
	}

	data := struct {
		Sections []NavItem
		Config   Config
	}{
		Sections: navItems,
		Config:   config,
	}

	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse index template: %w", err)
	}

	outputPath := filepath.Join(config.OutputDir, "index.html")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func generateSearchIndex(pages []Page, config Config) error {
	var searchPages []SearchPage
	re := regexp.MustCompile(`<[^>]*>`)

	for _, page := range pages {
		plainContent := re.ReplaceAllString(string(page.HTMLContent), " ")
		plainContent = strings.Join(strings.Fields(plainContent), " ")
		if len(plainContent) > 500 {
			plainContent = plainContent[:500] + "..."
		}

		searchPages = append(searchPages, SearchPage{
			Title:   page.Title,
			URL:     page.URL,
			Section: page.Section,
			Content: plainContent,
		})
	}

	index := SearchIndex{Pages: searchPages}
	jsonData, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	outputPath := filepath.Join(config.OutputDir, "search-index.json")
	return os.WriteFile(outputPath, jsonData, 0644)
}

func copyAssets(outputDir string) error {
	cssPath := filepath.Join(outputDir, "assets", "styles.css")
	if err := os.MkdirAll(filepath.Dir(cssPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(cssPath, []byte(cssStyles), 0644); err != nil {
		return err
	}

	jsPath := filepath.Join(outputDir, "assets", "app.js")
	if err := os.WriteFile(jsPath, []byte(jsScript), 0644); err != nil {
		return err
	}

	return nil
}

// Serve starts a local development server. dir/port fall back to the same
// defaults their CLI flags declare (dist/3000) — the `default:"..."` DTO
// struct tag isn't reliably applied when a flag is omitted entirely, so
// this defends the same way Build already does for its own Config fields.
func Serve(dir, port string) error {
	if dir == "" {
		dir = "dist"
	}
	if port == "" {
		port = "3000"
	}

	fmt.Printf("✉️  Starting development server at http://localhost:%s\n", port)
	fmt.Println("   Press Ctrl+C to stop")

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	return http.ListenAndServe(":"+port, nil)
}
