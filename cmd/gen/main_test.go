package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateRemovesStaleFiles(t *testing.T) {
	outDir := t.TempDir()

	// Pre-create a stale file that should not survive a fresh generate.
	staleDir := filepath.Join(outDir, "stale")
	if err := os.MkdirAll(staleDir, 0o755); err != nil {
		t.Fatalf("create stale dir: %v", err)
	}
	stalePath := filepath.Join(staleDir, "index.html")
	if err := os.WriteFile(stalePath, []byte("<html>stale</html>"), 0o644); err != nil {
		t.Fatalf("write stale file: %v", err)
	}

	if err := generate(outDir); err != nil {
		t.Fatalf("generate(%q): %v", outDir, err)
	}

	if _, err := os.Stat(stalePath); !os.IsNotExist(err) {
		t.Errorf("stale file %s should have been removed, but still exists", stalePath)
	}
}

func TestGenerate(t *testing.T) {
	outDir := t.TempDir()

	if err := generate(outDir); err != nil {
		t.Fatalf("generate(%q): %v", outDir, err)
	}

	// Every page that cmd/gen must produce.
	pages := []struct {
		rel   string // path relative to outDir
		title string // substring expected inside <title>
	}{
		{"index.html", "Home"},
		{"about/index.html", "About"},
		{"booking/index.html", "Booking"},
		{"projects/index.html", "Projects"},
		{"brand/index.html", "Brand"},
		{"partials/menu/index.html", "menu-btn"},             // HTMX fragment, contains the button id
		{"partials/menu-close/index.html", "menu-btn"},       // HTMX fragment, contains the button id
	}

	for _, p := range pages {
		path := filepath.Join(outDir, p.rel)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("expected file %s to exist: %v", p.rel, err)
			continue
		}
		html := string(data)

		if !strings.Contains(html, p.title) {
			t.Errorf("%s: expected to contain %q", p.rel, p.title)
		}
	}

	// Full pages (not partials) must reference the CSS stylesheet.
	fullPages := []string{
		"index.html",
		"about/index.html",
		"booking/index.html",
		"projects/index.html",
		"brand/index.html",
	}
	for _, rel := range fullPages {
		path := filepath.Join(outDir, rel)
		data, err := os.ReadFile(path)
		if err != nil {
			continue // already reported above
		}
		if !strings.Contains(string(data), "/static/css/app.css") {
			t.Errorf("%s: expected to contain /static/css/app.css", rel)
		}
	}

	// Static assets must be copied.
	// Check a known file exists in the output.
	jsPath := filepath.Join(outDir, "static", "js", "htmx.min.js")
	if _, err := os.Stat(jsPath); os.IsNotExist(err) {
		t.Errorf("expected static/js/htmx.min.js to be copied to output dir")
	}
}
