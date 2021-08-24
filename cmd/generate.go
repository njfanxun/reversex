package cmd

import (
	"fmt"

	"reversex/reverse"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewGenerateCommand() *cobra.Command {
	var filePath string
	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate a struct entity from database",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := reverse.NewReverseFromFile(filePath)
			if err != nil {
				pterm.Error.Printfln("%s", err.Error())
				return
			}
			err = r.Ping()
			if err != nil {
				pterm.Error.Printfln("%s", err.Error())
				return
			}
			tables, err := r.GetTableMetas()
			if err != nil {
				pterm.Error.Printfln("reversex get tables meta info error:%s", err.Error())
				return
			}
			err = r.PrepareRunReverse()
			if err != nil {
				pterm.Error.Printfln("reversex prepare to run error:%s", err.Error())
				return
			}
			p, _ := pterm.DefaultProgressbar.WithTotal(len(tables)).WithTitle(fmt.Sprintf("reverse %s file", r.GetLanguageName())).WithTitleStyle(pterm.NewStyle(pterm.FgGreen, pterm.Bold)).WithShowTitle(true).Start()
			for _, table := range tables {
				err = r.RunReverse(table)
				if err != nil {
					pterm.Error.Printfln("%s:%s", table.Name, err.Error())
				} else {
					pterm.Success.Printfln("generate file: %s ", table.Name)
				}
				p.Increment()
			}

		},
	}
	genCmd.Flags().StringVarP(&filePath, "file", "f", "", "yaml file to apply for reversex")
	_ = genCmd.MarkFlagRequired("file")
	return genCmd
}
