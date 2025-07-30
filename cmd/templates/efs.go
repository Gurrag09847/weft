package templates

import _ "embed"

//go:embed up.tmpl
var UpTemplate string

//go:embed down.tmpl
var DownTemplate string
