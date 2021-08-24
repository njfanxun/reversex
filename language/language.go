package language

import (
	"text/template"
)

type Language interface {
	GetName() string
	GetTemplate() string
	GetFunctions() template.FuncMap
	GetExtName() string
	FormatCodes(codes string) string
}
