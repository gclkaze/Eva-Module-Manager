package cmd

import (
	"emm/internal/app"

	"github.com/spf13/cobra"
)

func NewModuleUserGetCommand(application *app.EMMApp) *cobra.Command {
	var userModuleGetCmd = &cobra.Command{
		Use:     "mylist",
		Aliases: []string{"my"},
		Short:   "List my modules",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			tk, err := application.GetCurrentUserToken()
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}

			err = application.GetUserModules(tk)
			if err != nil {
				application.GetPrinter().Error(err)
			}
			return

		},
	}

	return userModuleGetCmd
}
