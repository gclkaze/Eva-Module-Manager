package cmd

import (
	"emm/internal/app"

	"github.com/spf13/cobra"
)

func NewUserParentCommand(application *app.EMMApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "User management-related commands.",
		Long:  "User management-related commands.",
	}

	cmd.AddCommand(
		NewUserListCommand(application),
		NewUserBanCommand(application),
		NewUserUnbanCommand(application),
	)
	return cmd
}
