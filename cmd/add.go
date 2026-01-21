package cmd

import "github.com/spf13/cobra"

var installedModule string = ""
var addCmd = &cobra.Command{
	Use:   "add module@version",
	Short: "Downloads and adds an EVA module to your EVA project.",
	Long:  "Downloads and adds an EVA module to your EVA project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
