package cmd

import (
	"emm/internal/app"
	"emm/internal/models/userinput"
	"emm/pkg/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func NewLoginCommand(application *app.EMMApp) *cobra.Command {
	var creds *userinput.LoginCreds

	var loginCmd = &cobra.Command{
		Use:     "login",
		Aliases: []string{"signin"},
		Short:   "User login to the Module Repository Server using an email and a password.",
		Long:    "User login to the Module Repository Server  using a email and a password, allowing him/her to perform Module management operations.",
		Args: func(cmd *cobra.Command, args []string) error {
			if !creds.AllInformationProvidedExceptPassword() {
				err := fmt.Errorf("in order to login, you will need to provide a valid registered email")
				application.GetPrinter().Error(err)
				application.SetOnError()
				return nil
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if application.IsOnError() {
				return
			}

			var err error
			pwd, err := utils.ReadPassword("Password: ")
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}
			//pwd := "mypass"
			creds.Password = pwd
			if creds.Password == "" {
				err = fmt.Errorf("no password provided")
				application.GetPrinter().Error(err)
				return
			}

			err = creds.AreValid()
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}

			if !creds.AllInformationProvided() {
				err = fmt.Errorf("in order to register, you will need to provide information such as your email & your password")
				application.GetPrinter().Error(err)
				return
			}
			err = application.UserLogin(creds)
			if err != nil {
				application.GetPrinter().Error(err)
			}
			return
		},
	}
	creds = userinput.NewLoginCreds()
	loginCmd.Flags().StringVarP(
		&creds.Email,
		"email",
		"e",
		"",
		"The user's email",
	)
	_ = loginCmd.MarkFlagRequired("email")
	return loginCmd
}
