package builder

const jsScript = `/**
 * Envelope Documentation - Interactive Features
 * Handles search, navigation, dark mode, and table of contents
 */

(function() {
  'use strict';

  // ========================================
  // Constants
  // ========================================
  const THEME_KEY = 'envelope-docs-theme';
  const SEARCH_INDEX_URL = 'search-index.json';

  // ========================================
  // State
  // ========================================
  let searchIndex = null;
  let searchResults = [];
  let selectedResultIndex = -1;

  // ========================================
  // DOM Elements
  // ========================================
  const elements = {
    themeSwitchBtns: document.querySelectorAll('.theme-switch-btn'),
    sidebar: document.querySelector('.sidebar'),
    sidebarOverlay: document.querySelector('.sidebar-overlay'),
    mobileMenuToggle: document.querySelector('.mobile-menu-toggle'),
    searchInput: document.querySelector('.search-input'),
    searchModal: document.querySelector('.search-modal'),
    searchModalInput: document.querySelector('.search-modal-input'),
    searchModalResults: document.querySelector('.search-modal-results'),
    navSections: document.querySelectorAll('.nav-section'),
    tocNav: document.querySelector('.toc-nav'),
    copyPageBtn: document.querySelector('.copy-page-button'),
    copyPageDropdownToggle: document.querySelector('.copy-page-dropdown-toggle'),
    copyPageDropdown: document.querySelector('.copy-page-dropdown')
  };

  // ========================================
  // Theme Management — three-way system/light/dark, matching the
  // sidebar's segmented control. "system" means no explicit choice: no
  // data-theme attribute is set at all, and prefers-color-scheme decides.
  // ========================================
  function initTheme() {
    const saved = localStorage.getItem(THEME_KEY); // 'light' | 'dark' | null
    applyTheme(saved);
  }

  function applyTheme(choice) {
    if (choice === 'light' || choice === 'dark') {
      document.documentElement.setAttribute('data-theme', choice);
    } else {
      document.documentElement.removeAttribute('data-theme');
      choice = 'system';
    }
    elements.themeSwitchBtns.forEach(btn => {
      btn.classList.toggle('active', btn.dataset.themeChoice === choice);
    });
  }

  function setThemeChoice(choice) {
    if (choice === 'system') {
      localStorage.removeItem(THEME_KEY);
      applyTheme(null);
    } else {
      localStorage.setItem(THEME_KEY, choice);
      applyTheme(choice);
    }
  }

  // ========================================
  // Mobile Navigation
  // ========================================
  function toggleSidebar() {
    elements.sidebar?.classList.toggle('open');
    elements.sidebarOverlay?.classList.toggle('active');
    document.body.style.overflow = elements.sidebar?.classList.contains('open') ? 'hidden' : '';
  }

  function closeSidebar() {
    elements.sidebar?.classList.remove('open');
    elements.sidebarOverlay?.classList.remove('active');
    document.body.style.overflow = '';
  }

  // ========================================
  // Navigation Sections
  // ========================================
  function initNavSections() {
    elements.navSections.forEach(section => {
      const header = section.querySelector('.nav-section-header');
      header?.addEventListener('click', () => {
        section.classList.toggle('active');
        const isExpanded = section.classList.contains('active');
        header.setAttribute('aria-expanded', isExpanded);
      });
    });
  }

  // ========================================
  // Search Functionality
  // ========================================
  async function loadSearchIndex() {
    try {
      // Get the base URL from any link on the page
      let baseUrl = '';
      const logoLink = document.querySelector('.logo');
      if (logoLink) {
        const href = logoLink.getAttribute('href');
        baseUrl = href.replace(/index\.html$/, '');
      }

      const response = await fetch(baseUrl + 'search-index.json');
      if (response.ok) {
        const data = await response.json();
        searchIndex = data.pages;
      }
    } catch (error) {
      console.warn('Could not load search index:', error);
    }
  }

  function openSearchModal() {
    elements.searchModal?.classList.add('active');
    elements.searchModalInput?.focus();
    document.body.style.overflow = 'hidden';
  }

  function closeSearchModal() {
    elements.searchModal?.classList.remove('active');
    elements.searchModalInput.value = '';
    elements.searchModalResults.innerHTML = '';
    searchResults = [];
    selectedResultIndex = -1;
    document.body.style.overflow = '';
  }

  function search(query) {
    if (!searchIndex || !query.trim()) {
      searchResults = [];
      renderSearchResults();
      return;
    }

    const terms = query.toLowerCase().split(/\s+/).filter(Boolean);

    searchResults = searchIndex
      .map(page => {
        let score = 0;
        const titleLower = page.title.toLowerCase();
        const contentLower = page.content.toLowerCase();

        terms.forEach(term => {
          // Title matches are weighted higher
          if (titleLower.includes(term)) {
            score += titleLower === term ? 100 : 50;
          }
          // Content matches
          if (contentLower.includes(term)) {
            score += 10;
          }
        });

        return { ...page, score };
      })
      .filter(page => page.score > 0)
      .sort((a, b) => b.score - a.score)
      .slice(0, 10);

    selectedResultIndex = searchResults.length > 0 ? 0 : -1;
    renderSearchResults();
  }

  function renderSearchResults() {
    if (!elements.searchModalResults) return;

    if (searchResults.length === 0) {
      const query = elements.searchModalInput?.value.trim();
      elements.searchModalResults.innerHTML = query
        ? '<div class="search-no-results">No results found</div>'
        : '';
      return;
    }

    elements.searchModalResults.innerHTML = searchResults.map((result, index) => {
      const sectionName = result.section === 'root' ? 'Overview' :
        result.section.split('-').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ');

      return ` + "`" + `
        <a href="${result.url}" class="search-result-item ${index === selectedResultIndex ? 'selected' : ''}" data-index="${index}">
          <div class="search-result-section">${sectionName}</div>
          <div class="search-result-title">${highlightMatch(result.title, elements.searchModalInput?.value || '')}</div>
          <div class="search-result-excerpt">${highlightMatch(result.content, elements.searchModalInput?.value || '')}</div>
        </a>
      ` + "`" + `;
    }).join('');
  }

  function highlightMatch(text, query) {
    if (!query.trim()) return escapeHtml(text);

    const terms = query.toLowerCase().split(/\s+/).filter(Boolean);
    let result = escapeHtml(text);

    terms.forEach(term => {
      const regex = new RegExp(` + "`" + `(${escapeRegex(term)})` + "`" + `, 'gi');
      result = result.replace(regex, '<mark>$1</mark>');
    });

    return result;
  }

  function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  function escapeRegex(string) {
    return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  }

  function navigateResults(direction) {
    if (searchResults.length === 0) return;

    if (direction === 'down') {
      selectedResultIndex = (selectedResultIndex + 1) % searchResults.length;
    } else {
      selectedResultIndex = (selectedResultIndex - 1 + searchResults.length) % searchResults.length;
    }

    renderSearchResults();

    // Scroll selected item into view
    const selectedItem = elements.searchModalResults?.querySelector('.selected');
    selectedItem?.scrollIntoView({ block: 'nearest' });
  }

  function selectResult() {
    if (selectedResultIndex >= 0 && selectedResultIndex < searchResults.length) {
      window.location.href = searchResults[selectedResultIndex].url;
    }
  }

  // ========================================
  // Table of Contents
  // ========================================
  function initTableOfContents() {
    if (!elements.tocNav) return;

    const headings = document.querySelectorAll('.prose h2, .prose h3');
    if (headings.length === 0) return;

    const tocItems = Array.from(headings).map(heading => {
      const id = heading.id || generateId(heading.textContent);
      heading.id = id;

      return {
        id,
        text: heading.textContent,
        level: heading.tagName.toLowerCase()
      };
    });

    elements.tocNav.innerHTML = tocItems.map(item => ` + "`" + `
      <a href="#${item.id}" class="toc-${item.level}">${escapeHtml(item.text)}</a>
    ` + "`" + `).join('');

    // Highlight active TOC item on scroll
    initTocHighlight(tocItems);
  }

  function generateId(text) {
    return text
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/(^-|-$)/g, '');
  }

  function initTocHighlight(tocItems) {
    const tocLinks = elements.tocNav?.querySelectorAll('a');
    if (!tocLinks || tocLinks.length === 0) return;

    const observer = new IntersectionObserver(entries => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          tocLinks.forEach(link => link.classList.remove('active'));
          const activeLink = elements.tocNav?.querySelector(` + "`" + `a[href="#${entry.target.id}"]` + "`" + `);
          activeLink?.classList.add('active');
        }
      });
    }, {
      rootMargin: '-20% 0px -80% 0px'
    });

    tocItems.forEach(item => {
      const element = document.getElementById(item.id);
      if (element) observer.observe(element);
    });
  }

  // ========================================
  // Event Listeners
  // ========================================
  function initEventListeners() {
    // Theme switch (system/light/dark)
    elements.themeSwitchBtns.forEach(btn => {
      btn.addEventListener('click', () => setThemeChoice(btn.dataset.themeChoice));
    });

    // Copy page (raw Markdown source) to clipboard
    elements.copyPageBtn?.addEventListener('click', async () => {
      const source = elements.copyPageBtn.dataset.source || '';
      try {
        await navigator.clipboard.writeText(source);
        const label = elements.copyPageBtn.querySelector('span');
        const original = label.textContent;
        elements.copyPageBtn.classList.add('copied');
        label.textContent = 'Copied';
        setTimeout(() => {
          elements.copyPageBtn.classList.remove('copied');
          label.textContent = original;
        }, 2000);
      } catch (err) {
        console.error('Failed to copy page:', err);
      }
    });

    // Copy page dropdown (View as Markdown, etc.)
    elements.copyPageDropdownToggle?.addEventListener('click', e => {
      e.stopPropagation();
      const isOpen = elements.copyPageDropdown?.classList.toggle('open');
      elements.copyPageDropdownToggle.setAttribute('aria-expanded', String(!!isOpen));
    });

    document.addEventListener('click', e => {
      if (elements.copyPageDropdown?.classList.contains('open') && !e.target.closest('.copy-page-group')) {
        elements.copyPageDropdown.classList.remove('open');
        elements.copyPageDropdownToggle?.setAttribute('aria-expanded', 'false');
      }
    });

    // Mobile menu
    elements.mobileMenuToggle?.addEventListener('click', toggleSidebar);
    elements.sidebarOverlay?.addEventListener('click', closeSidebar);

    // Search - sidebar input
    elements.searchInput?.addEventListener('focus', openSearchModal);

    // Search modal
    elements.searchModalInput?.addEventListener('input', e => {
      search(e.target.value);
    });

    // Search result hover
    elements.searchModalResults?.addEventListener('mouseover', e => {
      const item = e.target.closest('.search-result-item');
      if (item) {
        selectedResultIndex = parseInt(item.dataset.index, 10);
        renderSearchResults();
      }
    });

    // Close modal on click outside
    elements.searchModal?.addEventListener('click', e => {
      if (e.target === elements.searchModal) {
        closeSearchModal();
      }
    });

    // Keyboard shortcuts
    document.addEventListener('keydown', e => {
      // Open search with Cmd/Ctrl + K
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        openSearchModal();
      }

      // Close modal with Escape
      if (e.key === 'Escape' && elements.searchModal?.classList.contains('active')) {
        closeSearchModal();
      }

      // Close copy-page dropdown with Escape
      if (e.key === 'Escape' && elements.copyPageDropdown?.classList.contains('open')) {
        elements.copyPageDropdown.classList.remove('open');
        elements.copyPageDropdownToggle?.setAttribute('aria-expanded', 'false');
      }

      // Navigate results with arrow keys
      if (elements.searchModal?.classList.contains('active')) {
        if (e.key === 'ArrowDown') {
          e.preventDefault();
          navigateResults('down');
        } else if (e.key === 'ArrowUp') {
          e.preventDefault();
          navigateResults('up');
        } else if (e.key === 'Enter') {
          e.preventDefault();
          selectResult();
        }
      }
    });

    // Close sidebar on navigation
    document.querySelectorAll('.nav-link').forEach(link => {
      link.addEventListener('click', () => {
        if (window.innerWidth <= 1024) {
          closeSidebar();
        }
      });
    });
  }

  // ========================================
  // Copy Code Button
  // ========================================
  function initCodeCopyButtons() {
    document.querySelectorAll('.prose pre').forEach(pre => {
      if (pre.closest('.code-block')) return;

      const wrapper = document.createElement('div');
      wrapper.className = 'code-block';
      pre.parentNode.insertBefore(wrapper, pre);

      const bar = document.createElement('div');
      bar.className = 'code-block-bar';
      bar.innerHTML = ` + "`" + `
        <span class="code-block-dots"><span></span><span></span><span></span></span>
        <button class="copy-button" type="button" aria-label="Copy code">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
          </svg>
          <span>Copy</span>
        </button>
      ` + "`" + `;

      wrapper.appendChild(bar);
      wrapper.appendChild(pre);

      const button = bar.querySelector('.copy-button');
      button.addEventListener('click', async () => {
        const code = pre.querySelector('code')?.textContent || pre.textContent;
        try {
          await navigator.clipboard.writeText(code);
          button.classList.add('copied');
          button.querySelector('span').textContent = 'Copied';
          setTimeout(() => {
            button.classList.remove('copied');
            button.querySelector('span').textContent = 'Copy';
          }, 2000);
        } catch (err) {
          console.error('Failed to copy code:', err);
        }
      });
    });
  }

  // ========================================
  // Smooth Scroll for Anchor Links
  // ========================================
  function initSmoothScroll() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
      anchor.addEventListener('click', e => {
        const targetId = anchor.getAttribute('href').slice(1);
        const target = document.getElementById(targetId);
        if (target) {
          e.preventDefault();
          target.scrollIntoView({ behavior: 'smooth' });
          history.pushState(null, null, ` + "`" + `#${targetId}` + "`" + `);
        }
      });
    });
  }

  // ========================================
  // Initialize
  // ========================================
  function init() {
    initTheme();
    initNavSections();
    initTableOfContents();
    initEventListeners();
    initCodeCopyButtons();
    initSmoothScroll();
    loadSearchIndex();
  }

  // Run on DOM ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
`
