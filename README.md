# Reverse
A flexible and powerful tool for reverse database to model codes.<br />
Implementation reference [xorm.io/reverse](https://gitea.com/xorm/reverse) <br />
## Optimization

1. Optimized the speed of obtaining database table metadata.
2. Support special field definition of orm framework.
3. Support SqlType unsigned integer convert to uint for go language.
4. Support SqlType tinyint(1)  convert to bool for go language

## Installation
```go
go install github.com/njfanxun/reversex@latest
```
## Usage
```bash

reversex gen -f {dir}/config.yaml   # Generate a struct entity from database
reversex language                   # Print reversex supported language and framework
reversex version                    # Print the version of reversex

Use "reversex help [command]" for more information on a command.

```
## Configuration File
```yaml
kind: reverse           # --required, only kind reverse
name: mydb              # --required,this config name
source:
  database: mysql       # --required, database driver name (mysql|sqlite3|postgres|mssql)
  conn_str: ""          # --required, database DSN:[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...]
target:
  language: golang      # --required, development language(golang)
  include_tables:       # --optional, the include tables array list
  exclude_tables:       # --optional, the exclude tables array list
  mapper: snake         # --optional, table name to code class or struct mapping relationship(snake|gonic|same)
  table_prefix: ""      # --optional, struct name prefix
  orm_fields:           # --optional, xorm framework the special field identification, map[string]string
    created: created_at # key(orm field): value(database column)
    updated: updated_at
    version: version
  template_path: ""     # --optional, custom generate file template
  output_dir: ./models  # --required, output directory,the latest dir is package name

```

## Default goland template
```go
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
```
## Template Funcs
- UnTitle: Convert first charator of the word to lower.
- Upper: Convert word to all upper.
- TableMapper: Mapper method to convert table name to class/struct name.
- ColumnMapper: Mapper method to convert column name to class/struct field name.

## Golang Template Funcs
- Type: return column's golang type
- Tag: return golang struct tag for column
- Columns: return a table all columns to a array string
- Imports: return a go file all imports
- Package: return model files package 

## Dependencies
- github.com/cockroachdb/errors
- github.com/denisenkom/go-mssqldb
- github.com/go-sql-driver/mysql
- github.com/gobwas/glob
- github.com/lib/pq
- github.com/mattn/go-sqlite3
- github.com/novalagung/gubrak/v2
- github.com/pterm/pterm
- github.com/sirupsen/logrus
- github.com/spf13/cobra
- github.com/spf13/viper
- xorm.io/xorm 

## Todo
 support for languages such as java(MyBatis), node(Sequelize), etc.
