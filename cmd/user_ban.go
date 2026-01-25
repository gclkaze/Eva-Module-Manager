package cmd

import (
	"emm/internal/app"
	"emm/pkg/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func NewUserBanCommand(application *app.EMMApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ban <email>",
		Short: "Ban a user associated to a userID or email.",
		Long:  "Ban a user associated to a userID or email.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if utils.IsUint(args[0]) {
				banByID(cmd, args[0])
				return
			}
			banByEmail(cmd, args[0])
		},
	}
	return cmd
}

func banByID(cmd *cobra.Command, idValue string) {
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
	err = application.UserBanByID(ctx, tk, idUint)
	if err != nil {
		application.GetPrinter().Error(err)
		return
	}
}

func banByEmail(cmd *cobra.Command, email string) {
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
	err = application.UserBan(ctx, tk, email)
	if err != nil {
		application.GetPrinter().Error(err)
		return
	}
}
