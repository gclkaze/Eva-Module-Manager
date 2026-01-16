package cmd

import (
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current active user.",
	Long:  "Show current active user.",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := application.ShowCurrentUser()
		if err != nil {
			application.GetPrinter().Error(err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
