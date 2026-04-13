// Package ramirome exposes embedded static assets for the web server.
//
// The embed directive must live in a file co-located with the `static/`
// directory (module root) so the compiled binary ships as a single artifact
// with no external file dependencies.
package ramirome

import "embed"

// StaticFS contains everything under /static at build time.
//
//go:embed all:static
var StaticFS embed.FS
