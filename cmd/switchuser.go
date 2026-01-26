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
		Use:   "switchuser",
		Short: "Switch to another known and registered user.",
		Long:  "Switch to another known and registered user.",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if !utils.IsValidEmail(email) {
				err = fmt.Errorf("an invalid email was provided: '%s'. ", email)
				application.GetPrinter().Error(err)
				return nil
			}
			err = application.SwitchCurrentUser(email)
			if err != nil {
				application.GetPrinter().Error(err)
			}
			return nil
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
