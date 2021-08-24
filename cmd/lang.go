package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewLanguageCommand() *cobra.Command {

	langCmd := &cobra.Command{
		Use:   "language",
		Short: "Print reversex supported language and framework",
		Run: func(cmd *cobra.Command, args []string) {
			_ = pterm.DefaultTable.WithHasHeader().WithData(pterm.TableData{
				{"Language", "version", "framework"},
				{"goland", "1.16+", "xorm"},
			}).Render()

		},
	}

	return langCmd
}
