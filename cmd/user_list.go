package cmd

import (
	"emm/internal/app"

	"github.com/spf13/cobra"
)

func NewUserListCommand(application *app.EMMApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List users.",
		Long:    "List users.",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			tk, err := application.GetCurrentUserToken()
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}
			ctx := cmd.Context()
			err = application.UsersList(ctx, tk)
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}
		},
	}
	return cmd
}
