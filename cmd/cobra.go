package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const HelpTemplate = `{{TermTitleString .}}

{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}} 
`
const UsageTemplate = `{{ TermUsageString .}}`

const indent = " "

var (
	primaryStyle     = pterm.NewStyle(pterm.FgLightWhite)
	primaryBoldStyle = pterm.NewStyle(pterm.Bold, pterm.FgLightWhite)
)

func init() {
	cobra.AddTemplateFunc("TermTitleString", TermTitleString)
	cobra.AddTemplateFunc("TermUsageString", TermUsageString)
}

func TermTitleString(cmd *cobra.Command) string {
	return primaryBoldStyle.Sprintf("%s", cmd.Short)
}

func TermUsageString(cmd *cobra.Command) string {
	buf := bytes.NewBufferString(primaryBoldStyle.Sprintln("Usage:"))

	if cmd.HasAvailableSubCommands() {
		buf.WriteString(primaryStyle.Sprintln(LIndent(cmd.CommandPath()+" command [flag]", 2)))
	} else {
		buf.WriteString(primaryStyle.Sprintln(LIndent(cmd.CommandPath()+" [flag]", 2)))
	}
	if len(cmd.Commands()) > 0 {
		buf.WriteString("\n")
		buf.WriteString(primaryBoldStyle.Sprintln("Commands:"))
		for _, command := range cmd.Commands() {
			buf.WriteString(primaryStyle.Sprint(LIndent(RPad(command.Name(), command.NamePadding()), 2)))
			buf.WriteString(pterm.FgWhite.Sprintln(command.Short))
		}
	}
	if !cmd.HasAvailableSubCommands() {
		buf.WriteString("\n")
		buf.WriteString(primaryBoldStyle.Sprintln("Flags:"))
		buf.WriteString(cmd.Flags().FlagUsagesWrapped(pterm.GetTerminalWidth()))

	}
	if cmd.HasParent() {
		buf.WriteString(pterm.FgWhite.Sprintf("\nUse \"%s help [command]\" for more information on a command.", cmd.Parent().CommandPath()))
	} else {
		buf.WriteString(pterm.FgWhite.Sprintf("\nUse \"%s help [command]\" for more information on a command.", cmd.CommandPath()))
	}

	return buf.String()
}

func RPad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}
func LIndent(s string, count int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(indent, count), s)
}
