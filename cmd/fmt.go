package cmd

import (
	"github.com/gnagel/nginx-config-formatter/formatter"
	"github.com/spf13/cobra"
)

var formatterConfig * formatter.Fmt
var Indent string

// fmtCmd represents the fmt command
var fmtCmd = &cobra.Command{
	Use:   "fmt [--indent='\t'] [--backup] [--in-place] config_file",
	Short: "Nginx config file formatter",
	Long: `
# Nginx config file formatter

Format nginx configuration files in a standardized and consistent way:
* All lines are indented in uniform manner, with 4 spaces per level (default)
* Neighbouring empty lines are collapsed to at most two empty lines
* Curly braces placement follows Java convention
* Whitespaces are collapsed, except in comments an quotation marks
* Whitespaces in variable designators are removed: ${ my_variable } is collapsed to ${my_variable}
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := formatterConfig.Run()
		return err
	},
}


func init() {
	formatterConfig = &formatter.Fmt{}
	rootCmd.AddCommand(fmtCmd)
	fmtCmd.Flags().StringVarP(&formatterConfig.Indent, "indent", "i", "\t", "Set the indent amount (spaces or tab)")
	fmtCmd.Flags().BoolVarP(&formatterConfig.CreateBackup, "backup", "b", false, "Generate a backup")
	fmtCmd.Flags().BoolVarP(&formatterConfig.InPlace, "in-place", "w", false, "Write the config in-place")
	fmtCmd.Flags().BoolVarP(&formatterConfig.Verbose, "verbose", "v", false, "Print verbose messages")
	fmtCmd.Flags().StringVar(&formatterConfig.ConfigFile, "config", "c",  "Config file to process")
}
