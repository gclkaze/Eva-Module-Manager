package cmd

import (
	"emm/internal/app"
	"emm/pkg/utils"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewLogoutCommand(application *app.EMMApp) *cobra.Command {
	var email string
	var logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Logout the current active user or logout known user associated with the email input.",
		Long:  "Logout the current active user or logout known user associated with the email input.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			email = strings.TrimSpace(email)
			if email != "" {
				if !utils.IsValidEmail(email) {
					err = fmt.Errorf("an invalid email was provided: '%s'. ", email)
					application.GetPrinter().Error(err)
					return nil
				}
			}

			err = application.UserLogout(email)
			if err != nil {
				application.GetPrinter().Error(err)
			} else {
				application.GetPrinter().Info("Bye!")
			}
			return nil
		},
	}
	logoutCmd.Flags().StringVarP(
		&email,
		"email",
		"e",
		"",
		"The user's email",
	)
	return logoutCmd
}
