package reverse

import "github.com/spf13/viper"

type Config struct {
	Kind   string
	Name   string
	Source *Source
	Target *Target
}

type Source struct {
	Database string
	ConnStr  string
}

type Target struct {
	IncludeTables []string
	ExcludeTables []string
	Mapper        string
	TemplatePath  string
	OutputDir     string
	TablePrefix   string
	Language      string
	OrmFields     map[string]string
}

func NewConfig() *Config {
	return &Config{
		Kind: viper.GetString("kind"),
		Name: viper.GetString("name"),
		Source: &Source{
			Database: viper.GetString("source.database"),
			ConnStr:  viper.GetString("source.conn_str"),
		},
		Target: &Target{
			IncludeTables: viper.GetStringSlice("target.include_tables"),
			ExcludeTables: viper.GetStringSlice("target.exclude_tables"),
			Mapper:        viper.GetString("target.mapper"),
			TemplatePath:  viper.GetString("target.template_path"),
			OutputDir:     viper.GetString("target.output_dir"),
			TablePrefix:   viper.GetString("target.table_prefix"),
			Language:      viper.GetString("target.language"),
			OrmFields:     viper.GetStringMapString("target.orm_fields"),
		},
	}

}
