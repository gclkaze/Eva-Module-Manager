package cmd

import (
	"emm/internal/app"
	"emm/pkg/utils"

	"github.com/spf13/cobra"
)

func NewReleaseLowerCommand(application *app.EMMApp) *cobra.Command {
	var releaseLowerCmd = &cobra.Command{
		Use:   "lower module@version",
		Args:  cobra.ExactArgs(1),
		Short: "Change status of an accepted release of a Module with a specific Version to pending status",
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
			err = application.LowerRelease(tk, module, version)
			if err != nil {
				application.GetPrinter().Error(err)
			}
			return
		},
	}
	return releaseLowerCmd
}
