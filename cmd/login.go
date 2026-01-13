package cmd

import (
	"emm/internal/models/userinput"
	"fmt"

	"github.com/spf13/cobra"
)

var creds *userinput.LoginCreds

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "User login to the Module Repository Server using an email and a password.",
	Long:  "User login to the Module Repository Server  using a email and a password, allowing him/her to perform Module management operations",
	Args: func(cmd *cobra.Command, args []string) error {
		if !creds.AllInformationProvidedExceptPassword() {
			err := fmt.Errorf("in order to login, you will need to provide a valid registered email")
			application.GetPrinter().Error(err)
			return nil
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		/*pwd, err := utils.ReadPassword("Password: ")
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}*/
		pwd := "mypass"
		creds.Password = pwd
		if creds.Password == "" {
			err = fmt.Errorf("no password provided")
			application.GetPrinter().Error(err)
			return nil
		}

		/*	err = creds.AreValid()
			if err != nil {
				application.GetPrinter().Error(err)
				return nil
			}
		*/
		if !creds.AllInformationProvided() {
			err = fmt.Errorf("in order to register, you will need to provide information such as your email & your password")
			application.GetPrinter().Error(err)
			return nil
		}
		err = application.UserLogin(creds)
		if err != nil {
			application.GetPrinter().Error(err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	creds = userinput.NewLoginCreds()
	loginCmd.Flags().StringVarP(
		&creds.Email,
		"email",
		"u",
		"",
		"The user's email",
	)
	_ = registerCmd.MarkFlagRequired("email")
}
