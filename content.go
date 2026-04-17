// Package ramirome exposes embedded content assets for the markdown pipeline.
//
// The embed directive must live in a file co-located with the `content/`
// directory (module root) so the compiled binary ships with all markdown
// content baked in.
package ramirome

import "embed"

// ContentFS contains everything under /content at build time.
//
//go:embed all:content
var ContentFS embed.FS
