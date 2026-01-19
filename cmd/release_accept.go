package cmd

import (
	"emm/pkg/utils"

	"github.com/spf13/cobra"
)

var releaseAcceptCmd = &cobra.Command{
	Use:   "accept module@version",
	Args:  cobra.ExactArgs(1),
	Short: "Accept a suggested release of a Module with a specific Version",
	RunE: func(cmd *cobra.Command, args []string) error {
		module, version, err := utils.ParseModuleReleaseVersion(args[0])
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		tk, err := application.GetCurrentUserToken()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		err = application.AcceptRelease(tk, module, version)
		if err != nil {
			application.GetPrinter().Error(err)
		}
		return nil

	},
}

func init() {
	releaseCmd.AddCommand(releaseAcceptCmd)

}
