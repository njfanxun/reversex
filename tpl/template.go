package tpl

import (
	_ "embed"
	"strings"
	"text/template"

	"xorm.io/xorm/names"
)

//go:embed go.tpl
var DefaultGolangTemplate string

type FileTemplate struct {
	t      *template.Template
	funcs  template.FuncMap
	mapper string
}

func NewFileTemplate(mapper string) *FileTemplate {
	ft := &FileTemplate{
		t:      template.New("reverse"),
		funcs:  make(template.FuncMap),
		mapper: mapper,
	}
	// add template default function
	ft.funcs["UnTitle"] = UnTitle
	ft.funcs["Upper"] = Upper

	ft.funcs["TableMapper"] = GetMapperByName(ft.mapper).Table2Obj
	ft.funcs["ColumnMapper"] = GetMapperByName(ft.mapper).Table2Obj
	ft.t.Funcs(ft.funcs)
	return ft
}
func (ft *FileTemplate) FuncCount() int {
	return len(ft.funcs)
}
func (ft *FileTemplate) AddTemplateFunc(name string, tmplFunc interface{}) {
	ft.funcs[name] = tmplFunc
	ft.t.Funcs(ft.funcs)
}
func (ft *FileTemplate) AddTemplateFuncs(tmplFuncs template.FuncMap) {
	for k, v := range tmplFuncs {
		ft.funcs[k] = v
	}
	ft.t.Funcs(ft.funcs)
}
func (ft *FileTemplate) Parse(text string) (*template.Template, error) {
	return ft.t.Parse(text)
}

// ============ template default functions ========================================
/** UnTitle @Description: Convert first charator of the word to lower */
func UnTitle(src string) string {
	if src == "" {
		return ""
	}
	if len(src) == 1 {
		return strings.ToLower(string(src[0]))
	}
	return strings.ToLower(string(src[0])) + src[1:]
}

/** UpTitle @Description: Convert word to all upper */
func Upper(src string) string {
	if src == "" {
		return ""
	}

	return strings.ToUpper(src)
}

func GetMapperByName(mapname string) names.Mapper {
	switch mapname {
	case "gonic":
		return names.LintGonicMapper
	case "same":
		return names.SameMapper{}
	default:
		return names.SnakeMapper{}
	}
}
