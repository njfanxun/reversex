package reverse

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"reversex/language"
	"reversex/tpl"

	"github.com/cockroachdb/errors"
	"github.com/gobwas/glob"
	"github.com/novalagung/gubrak/v2"
	"github.com/spf13/viper"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

type Reverse struct {
	cfg    *Config
	lang   language.Language
	ft     *tpl.FileTemplate
	t      *template.Template
	engine *xorm.Engine
	ctx    context.Context
}

func viperBindFile(filePath string) error {
	dir := path.Dir(filePath)
	fileExt := strings.TrimPrefix(path.Ext(filePath), ".")
	fileName := strings.TrimSuffix(path.Base(filePath), path.Ext(filePath))
	viper.AddConfigPath(dir)
	viper.SetConfigType(fileExt)
	viper.SetConfigName(fileName)
	return viper.ReadInConfig()
}

func NewReverseFromFile(filePath string) (*Reverse, error) {
	err := viperBindFile(filePath)
	if err != nil {
		return nil, err
	}
	r := &Reverse{
		cfg:  NewConfig(),
		ctx:  context.Background(),
		lang: nil,
	}
	r.engine, err = xorm.NewEngine(r.cfg.Source.Database, r.cfg.Source.ConnStr)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// GetConfigFilePath /** @Description: 获取配置文件的路径 */
func (r *Reverse) GetConfigFilePath() string {
	return viper.ConfigFileUsed()
}

func (r *Reverse) Ping() error {
	return r.engine.DB().Ping()
}

func (r *Reverse) GetSource() *Source {
	return r.cfg.Source
}

func (r *Reverse) GetTarget() *Target {
	return r.cfg.Target
}
func (r *Reverse) GetConfigName() string {
	return r.cfg.Name
}
func (r *Reverse) GetLanguageName() string {
	if r.lang != nil {
		return r.lang.GetName()
	}
	return "unknown language"
}
func (r *Reverse) GetTableMetas() ([]*schemas.Table, error) {
	tables, err := r.engine.Dialect().GetTables(r.engine.DB(), r.ctx)
	if err != nil {
		return nil, err
	}
	// 按照配置过滤需要的tables
	r.filterTables(&tables)
	for _, table := range tables {
		colSeq, cols, err := r.engine.Dialect().GetColumns(r.engine.DB(), r.ctx, table.Name)
		if err != nil {
			return nil, err
		}
		for _, name := range colSeq {
			table.AddColumn(cols[name])
		}
		table.Indexes, err = r.engine.Dialect().GetIndexes(r.engine.DB(), r.ctx, table.Name)
		if err != nil {
			return nil, err
		}
		var seq int
		for _, index := range table.Indexes {
			for _, name := range index.Cols {
				parts := strings.Split(strings.TrimSpace(name), " ")
				if len(parts) > 1 {
					if parts[1] == "DESC" {
						seq = 1
					} else if parts[1] == "ASC" {
						seq = 0
					}
				}
				var colName = strings.Trim(parts[0], `"`)
				if col := table.GetColumn(colName); col != nil {
					col.Indexes[index.Name] = index.Type
				} else {
					return nil, errors.Errorf("Unknown col %s seq %d, in index %v of table %v, columns %v", name, seq, index.Name, table.Name, table.ColumnsSeq())
				}
			}
		}
	}
	return tables, nil
}

func (r *Reverse) PrepareRunReverse() error {
	switch r.cfg.Target.Language {
	case "golang":
		r.lang = &language.Golang{
			OrmFields:   r.cfg.Target.OrmFields,
			OutputDir:   r.cfg.Target.OutputDir,
			Mapper:      r.cfg.Target.Mapper,
			TablePrefix: r.cfg.Target.TablePrefix,
		}
	default:
		return errors.Errorf("unknown development language:%s", r.cfg.Target.Language)
	}

	r.ft = tpl.NewFileTemplate(r.cfg.Target.Mapper)

	// 添加各语言自己定义的模板方法
	r.ft.AddTemplateFuncs(r.lang.GetFunctions())

	// // load generate file template
	var bs []byte
	var err error
	if r.cfg.Target.TemplatePath != "" {
		bs, err = ioutil.ReadFile(r.cfg.Target.TemplatePath)
		if err != nil {
			return errors.Errorf("%s", err.Error())
		}
	} else {
		bs = []byte(r.lang.GetTemplate())
	}
	if bs == nil {
		return errors.New("you have to indicate template_path in yaml file for a language")
	}
	r.t, err = r.ft.Parse(string(bs))
	if err != nil {
		return err
	}
	err = os.MkdirAll(r.cfg.Target.OutputDir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (r *Reverse) RunReverse(table *schemas.Table) error {
	w, err := os.Create(filepath.Join(r.cfg.Target.OutputDir, table.Name+r.lang.GetExtName()))
	if err != nil {
		return err
	}
	defer func(w *os.File) {
		_ = w.Close()
	}(w)
	var buf = bytes.NewBufferString("")
	if r.cfg.Target.TablePrefix != "" {
		table.Name = strings.TrimPrefix(table.Name, r.cfg.Target.TablePrefix)
	}
	err = r.t.Execute(buf, map[string]interface{}{
		"Table": table,
	})
	if err != nil {
		return err
	}
	tplContent, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	var source string = r.lang.FormatCodes(string(tplContent))
	_, err = w.WriteString(source)
	return nil
}

func (r *Reverse) filterTables(tables *[]*schemas.Table) {
	chain := gubrak.From(*tables)
	// filter include tables

	if len(r.cfg.Target.IncludeTables) > 0 {
		chain.Filter(func(table *schemas.Table) bool {
			for _, includeTable := range r.cfg.Target.IncludeTables {
				g := glob.MustCompile(includeTable)
				if g.Match(table.Name) {
					return true
				}
			}
			return false
		})
	}
	//  filter exclude tables
	if len(r.cfg.Target.ExcludeTables) > 0 {
		chain.Filter(func(table *schemas.Table) bool {
			for _, excludeTable := range r.cfg.Target.ExcludeTables {
				g := glob.MustCompile(excludeTable)
				if g.Match(table.Name) {
					return false
				}
			}
			return true
		})
	}
	*tables = chain.Result().([]*schemas.Table)

}
