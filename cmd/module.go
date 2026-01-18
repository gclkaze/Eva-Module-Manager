package cmd

import "github.com/spf13/cobra"

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Module-related commands",
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}
