package cmd

import (
	"emm/internal/app"
	"emm/internal/models/userinput"

	"github.com/spf13/cobra"
)

func NewModuleUpdateCommand(application *app.EMMApp) *cobra.Command {
	var uploadedUpdatedCreds *userinput.UploadModuleUpdateParams

	var moduleUpdateCmd = &cobra.Command{
		Use:   "update [params...]",
		Short: "Update a module",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tk, err := application.GetCurrentUserToken()
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}

			err = uploadedUpdatedCreds.AllValid()
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}
			err = application.UpdateModule(tk, args, uploadedUpdatedCreds)
			if err != nil {
				application.GetPrinter().Error(err)
			}
		},
	}

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

	return moduleUpdateCmd
}
