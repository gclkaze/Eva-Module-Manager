package cmd

import (
	"emm/internal/app"
	"emm/pkg/utils"

	"github.com/spf13/cobra"
)

func NewReleaseRejectCommand(application *app.EMMApp) *cobra.Command {
	var releaseRejectCmd = &cobra.Command{
		Use:   "reject module@version",
		Args:  cobra.ExactArgs(1),
		Short: "Reject a suggested release of a Module with a specific Version",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := utils.IsValidSpecificModuleVersion(args[0])
			if !res {
				application.SetOnError()
				application.GetPrinter().Error(err)
				return
			}
			module, version, err := utils.ParseModuleReleaseVersion(args[0])
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}
			tk, err := application.GetCurrentUserToken()
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}
			err = application.RejectRelease(tk, module, version)
			if err != nil {
				application.GetPrinter().Error(err)
			}
		},
	}

	return releaseRejectCmd
}
