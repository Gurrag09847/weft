package templates

import _ "embed"

//go:embed up.tmpl
var UpTemplate string

//go:embed down.tmpl
var DownTemplate string

//go:embed handler.tmpl
var HandlerTemplate string

//go:embed service.tmpl
var ServiceTemplate string
