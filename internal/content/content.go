// Package content loads and renders Markdown files from an embedded filesystem
// at startup and exposes them as templ.Components for use in templates.
package content

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// cache stores rendered HTML keyed by path-without-extension
// (e.g. "about/intro" for content/about/intro.md).
var cache = make(map[string]string)

// md is the Goldmark instance with GFM extensions enabled.
var md = goldmark.New(goldmark.WithExtensions(extension.GFM, extension.Typographer))

// Init walks the embedded filesystem under "content/", reads every *.md file,
// renders it to HTML via Goldmark, and stores the result in the package-level
// cache. It must be called once at startup before any Get/MustGet calls.
func Init(fsys fs.FS) error {
	err := fs.WalkDir(fsys, "content", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		data, readErr := fs.ReadFile(fsys, path)
		if readErr != nil {
			return fmt.Errorf("reading %s: %w", path, readErr)
		}

		var buf bytes.Buffer
		if renderErr := md.Convert(data, &buf); renderErr != nil {
			return fmt.Errorf("rendering %s: %w", path, renderErr)
		}

		// Strip "content/" prefix and ".md" extension to form the key.
		key := strings.TrimPrefix(path, "content/")
		key = strings.TrimSuffix(key, ".md")

		cache[key] = buf.String()
		return nil
	})
	if err != nil {
		return fmt.Errorf("walking content dir: %w", err)
	}

	log.Printf("content: loaded %d markdown files", len(cache))
	return nil
}

// Get returns a templ.Component that renders the cached HTML for the given key.
// If the key is missing, it logs a warning and returns a component that renders
// an empty string. It does NOT panic — safe for use in production request paths.
func Get(key string) templ.Component {
	html, ok := cache[key]
	if !ok {
		log.Printf("content: warning: key %q not found", key)
		return templ.Raw("")
	}
	return templ.Raw(html)
}

// MustGet returns a templ.Component for the given key. It panics if the key
// is missing. Use this from templ files so missing content is caught loudly
// during development.
func MustGet(key string) templ.Component {
	html, ok := cache[key]
	if !ok {
		panic(fmt.Sprintf("content: key %q not found", key))
	}
	return templ.Raw(html)
}
