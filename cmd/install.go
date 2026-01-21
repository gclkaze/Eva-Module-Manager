package cmd

import "github.com/spf13/cobra"

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Reads and installs the EVA modules according to eva.json.",
	Long:  "Reads and installs the EVA modules according to eva.json.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
