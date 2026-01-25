package cmd

import (
	"emm/internal/app"
	"emm/pkg/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func NewUserUnbanCommand(application *app.EMMApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unban <email>",
		Short: "Unban an already banned user associated to a userID or an email.",
		Long:  "Unban an already banned user associated to a userID or an email.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if utils.IsUint(args[0]) {
				unbanByID(cmd, args[0])
				return
			}
			unbanByEmail(cmd, args[0])
		},
	}
	return cmd
}

func unbanByID(cmd *cobra.Command, idValue string) {
	idUint, err := utils.StringToUint(idValue)
	if err != nil {
		application.GetPrinter().Error(err)
		return
	}
	tk, err := application.GetCurrentUserToken()
	if err != nil {
		application.GetPrinter().Error(err)
		return
	}
	ctx := cmd.Context()
	err = application.UserUnbanByID(ctx, tk, idUint)
	if err != nil {
		application.GetPrinter().Error(err)
		return
	}
}

func unbanByEmail(cmd *cobra.Command, email string) {
	if !utils.IsValidEmail(email) {
		application.GetPrinter().Error(fmt.Errorf("invalid email provided"))
		return
	}
	tk, err := application.GetCurrentUserToken()
	if err != nil {
		application.GetPrinter().Error(err)
		return
	}
	ctx := cmd.Context()
	err = application.UserUnban(ctx, tk, email)
	if err != nil {
		application.GetPrinter().Error(err)
		return
	}
}
