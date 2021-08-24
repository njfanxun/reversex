package {{Package}}

{{$ips := Imports .Table}}
{{$ilen := len $ips}}
{{if gt $ilen 0}}
    import (
    {{range $ips}}"{{.}}"{{end}}
    )
{{end}}

{{with .Table}}
    type {{TableMapper .Name}} struct {
    {{$table := .}}
    {{range .ColumnsSeq}}{{$col := $table.GetColumn .}}    {{ColumnMapper $col.Name}}    {{Type $col}} `json:"{{$col.Name}}" {{Tag $table $col}}`
    {{end}}
    }

    func ({{TableMapper .Name}}) TableName() string {
    return "{{$table.Name}}"
    }

    func ({{TableMapper .Name}}) Cols() []string {

    return {{ Columns .}}
    }
{{end}}
