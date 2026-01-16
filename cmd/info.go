package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var module string

var showCmd = &cobra.Command{
	Use:   "info",
	Short: "Show module or module release information",
	Long:  "Show module information or module release information",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			application.SetOnError()
			return fmt.Errorf("provide the module name or module-name@version for module/release information")
		}
		module = args[0]
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if application.IsOnError() {
			return nil
		}
		err := application.GetModuleInfo(module)
		return err
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
