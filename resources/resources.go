package resources

import "embed"

//go:embed views/*.html
var ViewsFS embed.FS
