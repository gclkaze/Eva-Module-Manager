package cmd

import (
	"emm/internal/app"
	"emm/pkg/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func NewSwitchUserCommand(application *app.EMMApp) *cobra.Command {
	var email string

	var switchCmd = &cobra.Command{
		Use:     "switchuser",
		Aliases: []string{"sw"},
		Short:   "Switch to another known and registered user.",
		Long:    "Switch to another known and registered user.",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if !utils.IsValidEmail(email) {
				err = fmt.Errorf("an invalid email was provided: '%s'. ", email)
				application.GetPrinter().Error(err)
				return
			}
			err = application.SwitchCurrentUser(email)
			if err != nil {
				application.GetPrinter().Error(err)
			}
		},
	}
	switchCmd.Flags().StringVarP(
		&email,
		"email",
		"e",
		"",
		"The user's email",
	)
	_ = switchCmd.MarkFlagRequired("email")

	return switchCmd
}
