package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userModuleGetCmd = &cobra.Command{
	Use:   "mylist [params...]",
	Short: "List my modules",
	RunE: func(cmd *cobra.Command, args []string) error {
		tk, err := application.GetCurrentUserToken()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}

		err = application.GetUserModules(tk)
		if err != nil {
			application.GetPrinter().Error(err)
		}
		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Upload called with params:")
		for i, arg := range args {
			fmt.Printf("  %d: %s\n", i+1, arg)
		}
	},
}

func init() {
	moduleCmd.AddCommand(userModuleGetCmd)

}
