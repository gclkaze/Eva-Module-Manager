package cmd

import (
	"github.com/spf13/cobra"
)

var registrationEmail string
var registrationPwd string

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register to the Module Repository Server using an email and a password.",
	Long:  "Register to the Module Repository Server using a username and a password, allowing him/her to perform Module management operations",
	Args: func(cmd *cobra.Command, args []string) error {
		/*		if len(args) != 1 {
					return fmt.Errorf("provide the module name or module-name@version for module/release information")
				}
				module = args[0]*/
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		//err := application.GetModuleInfo(module)
		return nil //err
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringVarP(
		&registrationEmail,
		"email",
		"u",
		"",
		"The user's email",
	)

	registerCmd.Flags().StringVarP(
		&registrationPwd,
		"password",
		"p",
		"",
		"The user's password",
	)
}
