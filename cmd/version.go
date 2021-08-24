package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const VERSION = "1.0.0"

func NewVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of reversex.",
		Run: func(cmd *cobra.Command, args []string) {
			pterm.FgLightCyan.Printfln("reversex version: %s", VERSION)
		},
	}
	return versionCmd
}
