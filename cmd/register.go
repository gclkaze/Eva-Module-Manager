package cmd

import (
	"emm/internal/models/userinput"
	"emm/pkg/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var registrationCreds *userinput.RegistrationCreds

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register to the Module Repository Server using an email, a password, a handle and your first name and last name.",
	Long:  "Register to the Module Repository Server using a username and a password, a handle and your first name and last name, in order to perform Module management operations and contribute to the EVA Module Developer.",
	Args: func(cmd *cobra.Command, args []string) error {
		if !registrationCreds.AllInformationProvidedExceptPassword() {
			err := fmt.Errorf("in order to register, you will need to provide information such as your email, first & last name, a password and a handle; a unique identifier for your profile")
			application.GetPrinter().Error(err)
			return nil
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		pwd, err := utils.ReadPassword("Password: ")
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		registrationCreds.Password = pwd
		if registrationCreds.Password == "" {
			err = fmt.Errorf("no password provided")
			application.GetPrinter().Error(err)
			return nil
		}

		err = registrationCreds.AreValid()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}

		if !registrationCreds.AllInformationProvided() {
			err = fmt.Errorf("in order to register, you will need to provide information such as your email, first & last name, a password and a handle; a unique identifier for your profile")
			application.GetPrinter().Error(err)
			return nil
		}
		err = application.UserRegister(registrationCreds)
		if err != nil {
			application.GetPrinter().Error(err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registrationCreds = userinput.NewRegistrationCreds()
	registerCmd.Flags().StringVarP(
		&registrationCreds.Email,
		"email",
		"u",
		"",
		"The user's email",
	)

	registerCmd.Flags().StringVarP(
		&registrationCreds.FirstName,
		"firstname",
		"f",
		"",
		"The user's first name",
	)

	registerCmd.Flags().StringVarP(
		&registrationCreds.LastName,
		"lastname",
		"l",
		"",
		"The user's last name",
	)

	registerCmd.Flags().StringVarP(
		&registrationCreds.Handle,
		"handle",
		"a",
		"",
		"The user's handle",
	)

	_ = registerCmd.MarkFlagRequired("email")
	_ = registerCmd.MarkFlagRequired("firstname")
	_ = registerCmd.MarkFlagRequired("lastname")
	_ = registerCmd.MarkFlagRequired("handle")
}
