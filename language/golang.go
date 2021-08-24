package language

import (
	"bytes"
	"fmt"
	"go/format"
	"path"
	"sort"
	"strings"
	"text/template"

	"reversex/tpl"

	"xorm.io/xorm/schemas"
)

type Golang struct {
	OrmFields   map[string]string
	OutputDir   string
	Mapper      string
	TablePrefix string
}

// =========================== implement language interface ==================================================

func (g *Golang) GetName() string {
	return "golang"
}

func (g *Golang) GetExtName() string {
	return ".go"
}

func (g *Golang) GetTemplate() string {
	return tpl.DefaultGolangTemplate
}

func (g *Golang) GetFunctions() template.FuncMap {
	return template.FuncMap{
		"Type":    GetTypeString,
		"Columns": GetColumns,
		"Imports": GetImporters,
		"Tag":     g.Tag,
		"Package": g.GetPackage,
	}
}

func (g *Golang) FormatCodes(codes string) string {
	source, err := format.Source([]byte(codes))
	if err != nil {
		return codes
	}
	return string(source)
}

// =========================== custom go template functions ==================================================

func GetColumns(table *schemas.Table) string {
	var buf = bytes.NewBufferString("[]string{")
	for _, s := range table.ColumnsSeq() {
		buf.WriteString(fmt.Sprintf("\"%s\",", s))
	}
	buf.Truncate(buf.Len() - 1)
	buf.WriteString("}")
	return buf.String()
}

func GetTypeString(col *schemas.Column) string {
	st := col.SQLType
	t := SQLType2GoType(st)
	s := t.String()
	if s == "[]uint8" {
		return "[]byte"
	}
	return s
}

func GetImporters(table *schemas.Table) []string {

	var results []string = make([]string, 0)
	for _, column := range table.Columns() {
		if GetTypeString(column) == "time.Time" {
			results = append(results, "time")
			break
		}
	}
	return results

}

func (g *Golang) Tag(table *schemas.Table, col *schemas.Column) string {
	var res []string
	if col.IsPrimaryKey {
		res = append(res, "pk")
	}
	if col.IsAutoIncrement {
		res = append(res, "autoincr")
	}
	if _, ok := g.OrmFields[col.Name]; ok {
		if col.SQLType.IsTime() {
			res = append(res, col.Name)
		}
	}
	if !col.Nullable {
		res = append(res, "not null")
	}
	if col.Default != "" {
		res = append(res, "default "+col.Default)
	}
	names := make([]string, 0, len(col.Indexes))
	for name := range col.Indexes {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		index := table.Indexes[name]
		var uiStr string
		if index.Type == schemas.UniqueType {
			uiStr = "unique"
		} else if index.Type == schemas.IndexType {
			uiStr = "index"
		}
		if len(index.Cols) > 1 {
			uiStr += "(" + index.Name + ")"
		}
		res = append(res, uiStr)
	}
	nStr := col.SQLType.Name
	if col.Length != 0 {
		if col.Length2 != 0 {
			nStr += fmt.Sprintf("(%v,%v)", col.Length, col.Length2)
		} else {
			nStr += fmt.Sprintf("(%v)", col.Length)
		}
	} else if len(col.EnumOptions) > 0 { // enum
		nStr += "("
		opts := ""

		enumOptions := make([]string, 0, len(col.EnumOptions))
		for enumOption := range col.EnumOptions {
			enumOptions = append(enumOptions, enumOption)
		}
		sort.Strings(enumOptions)

		for _, v := range enumOptions {
			opts += fmt.Sprintf(",'%v'", v)
		}
		nStr += strings.TrimLeft(opts, ",")
		nStr += ")"
	} else if len(col.SetOptions) > 0 {
		nStr += "("
		opts := ""

		setOptions := make([]string, 0, len(col.SetOptions))
		for setOption := range col.SetOptions {
			setOptions = append(setOptions, setOption)
		}
		sort.Strings(setOptions)

		for _, v := range setOptions {
			opts += fmt.Sprintf(",'%v'", v)
		}
		nStr += strings.TrimLeft(opts, ",")
		nStr += ")"
	}
	res = append(res, nStr)
	if len(res) > 0 {
		return fmt.Sprintf(`xorm:"'%s' %s"`, col.Name, strings.Join(res, " "))

	}
	return ""
}

func (g *Golang) GetPackage() string {
	pkg := path.Base(g.OutputDir)
	if pkg == "" {
		return "models"
	}
	return pkg
}
