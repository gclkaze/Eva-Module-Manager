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
		if len(args) != 1 {
			err := fmt.Errorf("both email and password are required")
			application.GetPrinter().Error(err)
			return nil
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		//err := application.GetModuleInfo(module)
		return nil //err
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	creds = userinput.NewLoginCreds()
	loginCmd.Flags().StringVarP(
		&creds.Email,
		"username",
		"u",
		"",
		"The user's email",
	)

	loginCmd.Flags().StringVarP(
		&creds.Password,
		"password",
		"p",
		"",
		"The user's password",
	)
}
