package cmd

import (
	"emm/internal/app"

	"github.com/spf13/cobra"
)

func NewModuleUserGetCommand(application *app.EMMApp) *cobra.Command {
	var userModuleGetCmd = &cobra.Command{
		Use:   "mylist",
		Short: "List my modules",
		Args:  cobra.ExactArgs(0),
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
	}

	return userModuleGetCmd
}
