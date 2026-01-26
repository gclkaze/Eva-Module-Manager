package cmd

import (
	"emm/internal/app"

	"github.com/spf13/cobra"
)

func NewWhoamiCommand(application *app.EMMApp) *cobra.Command {
	var whoamiCmd = &cobra.Command{
		Use:   "whoami",
		Short: "Show current active user.",
		Long:  "Show current active user.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := application.ShowCurrentUser()
			if err != nil {
				application.GetPrinter().Error(err)
			}
			return nil
		},
	}
	return whoamiCmd
}
