// Package main is a static-site generator that renders all templ pages
// into ./public/ as plain HTML files suitable for static hosting.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/a-h/templ"

	ramirome "github.com/rejkpp/ramiro.me"
	"github.com/rejkpp/ramiro.me/internal/content"
	"github.com/rejkpp/ramiro.me/templates/layout"
	"github.com/rejkpp/ramiro.me/templates/pages"
)

// page defines a route to render as a full HTML page (wrapped in layout.Base).
type page struct {
	route     string          // URL path, e.g. "/" or "/about"
	component templ.Component // the templ page component
}

// partial defines a route to render as an HTMX swap fragment (no layout).
type partial struct {
	route     string
	component templ.Component
}

func main() {
	if err := generate("public"); err != nil {
		log.Fatalf("generate: %v", err)
	}
	log.Println("static site generated in ./public/")
}

// generate renders all pages and partials to outDir, and copies static assets.
func generate(outDir string) error {
	// Initialize markdown content (required by Home, About pages).
	if err := content.Init(ramirome.ContentFS); err != nil {
		return fmt.Errorf("content init: %w", err)
	}

	allPages := []page{
		{"/", pages.Home()},
		{"/about", pages.About()},
		{"/booking", pages.Booking()},
		{"/projects", pages.Projects()},
		{"/brand", pages.Brand()},
	}

	allPartials := []partial{
		{"/partials/menu", layout.Menu()},
		{"/partials/menu-close", layout.MenuClose()},
	}

	ctx := context.Background()

	for _, p := range allPages {
		if err := renderToFile(ctx, outDir, p.route, p.component); err != nil {
			return fmt.Errorf("render page %s: %w", p.route, err)
		}
	}

	for _, p := range allPartials {
		if err := renderToFile(ctx, outDir, p.route, p.component); err != nil {
			return fmt.Errorf("render partial %s: %w", p.route, err)
		}
	}

	if err := copyStaticAssets(outDir); err != nil {
		return fmt.Errorf("copy static: %w", err)
	}

	return nil
}

// renderToFile renders a templ component to outDir/<route>/index.html.
// The root route "/" maps to outDir/index.html.
func renderToFile(ctx context.Context, outDir, route string, c templ.Component) error {
	var rel string
	if route == "/" {
		rel = "index.html"
	} else {
		rel = filepath.Join(route, "index.html")
	}

	dest := filepath.Join(outDir, rel)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := c.Render(ctx, &buf); err != nil {
		return fmt.Errorf("templ render: %w", err)
	}

	return os.WriteFile(dest, buf.Bytes(), 0o644)
}

// copyStaticAssets walks the embedded static/ directory and copies every file
// into outDir/static/.
func copyStaticAssets(outDir string) error {
	staticFS, err := fs.Sub(ramirome.StaticFS, "static")
	if err != nil {
		return fmt.Errorf("sub static fs: %w", err)
	}

	return fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		dest := filepath.Join(outDir, "static", path)

		if d.IsDir() {
			return os.MkdirAll(dest, 0o755)
		}

		src, openErr := staticFS.Open(path)
		if openErr != nil {
			return fmt.Errorf("open %s: %w", path, openErr)
		}
		defer src.Close()

		if mkErr := os.MkdirAll(filepath.Dir(dest), 0o755); mkErr != nil {
			return mkErr
		}

		out, createErr := os.Create(dest)
		if createErr != nil {
			return fmt.Errorf("create %s: %w", dest, createErr)
		}
		defer out.Close()

		if _, cpErr := io.Copy(out, src); cpErr != nil {
			return fmt.Errorf("copy %s: %w", path, cpErr)
		}

		return nil
	})
}
