// Package web embeds templates and static assets so the compiled
// binary is fully self-contained in production.
package web

import "embed"

//go:embed templates
var Templates embed.FS

//go:embed static
var Static embed.FS
