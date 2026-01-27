package cmd

import (
	"emm/internal/app"

	"github.com/spf13/cobra"
)

func NewWhoamiCommand(application *app.EMMApp) *cobra.Command {
	var whoamiCmd = &cobra.Command{
		Use:     "whoami",
		Aliases: []string{"w"},
		Short:   "Show current active user.",
		Long:    "Show current active user.",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := application.ShowCurrentUser()
			if err != nil {
				application.GetPrinter().Error(err)
			}
		},
	}
	return whoamiCmd
}
