package cmd

import (
	"emm/internal/models/userinput"

	"github.com/spf13/cobra"
)

var uploadedUpdatedCreds *userinput.UploadModuleUpdateParams

var moduleUpdateCmd = &cobra.Command{
	Use:   "update [params...]",
	Short: "Update a module",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tk, err := application.GetCurrentUserToken()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}

		err = uploadedUpdatedCreds.AllValid()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		err = application.UpdateModule(tk, args, uploadedUpdatedCreds)
		if err != nil {
			application.GetPrinter().Error(err)
		}
		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {
		/*		fmt.Println("Update called with params:")
				for i, arg := range args {
					fmt.Printf("  %d: %s\n", i+1, arg)
				}*/
	},
}

func init() {
	uploadedUpdatedCreds = userinput.NewUploadModuleUpdateParams()

	moduleUpdateCmd.Flags().StringVarP(
		&uploadedUpdatedCreds.ModuleRepr,
		"module-name",
		"m",
		"",
		"module name (required)",
	)

	moduleUpdateCmd.Flags().StringVarP(
		&uploadedUpdatedCreds.Title,
		"title",
		"n",
		"",
		"Module title (required)",
	)

	moduleUpdateCmd.Flags().StringVarP(
		&uploadedUpdatedCreds.Repr,
		"repr",
		"r",
		"",
		"Module identifier (required)",
	)

	moduleUpdateCmd.Flags().StringVarP(
		&uploadedUpdatedCreds.Tags,
		"tags",
		"t",
		"",
		"Comma-separated tags",
	)

	moduleUpdateCmd.Flags().StringVarP(
		&uploadedUpdatedCreds.Description,
		"description",
		"d",
		"",
		"Module description",
	)

	_ = moduleUpdateCmd.MarkFlagRequired("title")
	_ = moduleUpdateCmd.MarkFlagRequired("repr")
	_ = moduleUpdateCmd.MarkFlagRequired("module-name")

	moduleCmd.AddCommand(moduleUpdateCmd)

}
