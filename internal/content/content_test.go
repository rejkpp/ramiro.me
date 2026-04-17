package content

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"testing/fstest"
)

func TestInit_PopulatesCache(t *testing.T) {
	fs := fstest.MapFS{
		"content/about/intro.md": &fstest.MapFile{
			Data: []byte("# Hello\n\nWorld paragraph."),
		},
		"content/about/closing.md": &fstest.MapFile{
			Data: []byte("Goodbye **world**."),
		},
	}

	resetCache()
	if err := Init(fs); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	if len(cache) != 2 {
		t.Fatalf("expected 2 cache entries, got %d", len(cache))
	}

	if _, ok := cache["about/intro"]; !ok {
		t.Error("missing cache key 'about/intro'")
	}
	if _, ok := cache["about/closing"]; !ok {
		t.Error("missing cache key 'about/closing'")
	}
}

func TestInit_RendersMarkdown(t *testing.T) {
	fs := fstest.MapFS{
		"content/page.md": &fstest.MapFile{
			Data: []byte("# hello"),
		},
	}

	resetCache()
	if err := Init(fs); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	html := cache["page"]
	if !strings.Contains(html, "<h1>hello</h1>") {
		t.Errorf("expected rendered <h1>, got: %s", html)
	}
}

func TestInit_GFMExtension(t *testing.T) {
	fs := fstest.MapFS{
		"content/gfm.md": &fstest.MapFile{
			Data: []byte("~~strikethrough~~"),
		},
	}

	resetCache()
	if err := Init(fs); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	html := cache["gfm"]
	if !strings.Contains(html, "<del>strikethrough</del>") {
		t.Errorf("expected GFM strikethrough, got: %s", html)
	}
}

func TestGet_UnknownKey_ReturnsEmptyNoPanic(t *testing.T) {
	resetCache()
	comp := Get("nonexistent/key")
	if comp == nil {
		t.Fatal("Get returned nil for unknown key")
	}

	var buf bytes.Buffer
	err := comp.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if buf.String() != "" {
		t.Errorf("expected empty string for unknown key, got: %q", buf.String())
	}
}

func TestGet_KnownKey_ReturnsRenderedHTML(t *testing.T) {
	fs := fstest.MapFS{
		"content/test.md": &fstest.MapFile{
			Data: []byte("**bold text**"),
		},
	}

	resetCache()
	if err := Init(fs); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	comp := Get("test")
	var buf bytes.Buffer
	err := comp.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "<strong>bold text</strong>") {
		t.Errorf("expected bold HTML, got: %s", output)
	}
}

func TestMustGet_UnknownKey_Panics(t *testing.T) {
	resetCache()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("MustGet did not panic for unknown key")
		}
	}()

	MustGet("nonexistent/key")
}

func TestMustGet_KnownKey_ReturnsComponent(t *testing.T) {
	fs := fstest.MapFS{
		"content/known.md": &fstest.MapFile{
			Data: []byte("known content"),
		},
	}

	resetCache()
	if err := Init(fs); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	comp := MustGet("known")
	var buf bytes.Buffer
	err := comp.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if !strings.Contains(buf.String(), "known content") {
		t.Errorf("expected 'known content', got: %s", buf.String())
	}
}

func TestInit_IgnoresNonMarkdownFiles(t *testing.T) {
	fs := fstest.MapFS{
		"content/readme.txt": &fstest.MapFile{
			Data: []byte("not markdown"),
		},
		"content/page.md": &fstest.MapFile{
			Data: []byte("markdown"),
		},
	}

	resetCache()
	if err := Init(fs); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	if len(cache) != 1 {
		t.Fatalf("expected 1 cache entry (only .md files), got %d", len(cache))
	}
}

func TestInit_EmptyFS(t *testing.T) {
	fs := fstest.MapFS{
		"content/.gitkeep": &fstest.MapFile{
			Data: []byte(""),
		},
	}

	resetCache()
	err := Init(fs)
	if err != nil {
		t.Fatalf("Init should not error on empty content dir: %v", err)
	}

	if len(cache) != 0 {
		t.Fatalf("expected 0 cache entries, got %d", len(cache))
	}
}

func TestInit_NestedDirectories(t *testing.T) {
	fs := fstest.MapFS{
		"content/a/b/deep.md": &fstest.MapFile{
			Data: []byte("deep content"),
		},
	}

	resetCache()
	if err := Init(fs); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	if _, ok := cache["a/b/deep"]; !ok {
		t.Errorf("expected cache key 'a/b/deep', keys: %v", cacheKeys())
	}
}

// resetCache clears the package-level cache for test isolation.
func resetCache() {
	cache = make(map[string]string)
}

// cacheKeys returns all keys in the cache (for debugging).
func cacheKeys() []string {
	keys := make([]string, 0, len(cache))
	for k := range cache {
		keys = append(keys, k)
	}
	return keys
}
