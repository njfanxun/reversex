package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "reversex",
	Short:         "A flexible and powerful tool for reverse database to model codes",
	SilenceErrors: true,
	SilenceUsage:  true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
	},
}

func init() {
	rootCmd.SetHelpTemplate(HelpTemplate)
	rootCmd.SetUsageTemplate(UsageTemplate)
	rootCmd.AddCommand(NewGenerateCommand())
	rootCmd.AddCommand(NewLanguageCommand())
	rootCmd.AddCommand(NewVersionCommand())
}

func ExecuteCommand() error {
	err := rootCmd.Execute()
	if err != nil {
		pterm.Error.Printfln("%s", err.Error())
	}
	return err
}
