package cmd

import (
	"emm/internal/app"
	"emm/internal/models/userinput"

	"github.com/spf13/cobra"
)

func NewModuleUploadCommand(application *app.EMMApp) *cobra.Command {
	var uploadedCreds *userinput.UploadParams

	var moduleUploadCmd = &cobra.Command{
		Use:   "upload [params...]",
		Short: "Upload a module",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tk, err := application.GetCurrentUserToken()
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}

			err = uploadedCreds.AllValid()
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}
			err = application.UploadModule(tk, args, uploadedCreds)
			if err != nil {
				application.GetPrinter().Error(err)
			}
			return

		},
	}
	uploadedCreds = userinput.NewUploadParams()
	moduleUploadCmd.Flags().StringVarP(
		&uploadedCreds.Title,
		"title",
		"n",
		"",
		"Module title (required)",
	)

	moduleUploadCmd.Flags().StringVarP(
		&uploadedCreds.Repr,
		"repr",
		"r",
		"",
		"Module identifier (required)",
	)

	moduleUploadCmd.Flags().StringVarP(
		&uploadedCreds.Tags,
		"tags",
		"t",
		"",
		"Comma-separated tags",
	)

	moduleUploadCmd.Flags().StringVarP(
		&uploadedCreds.Description,
		"description",
		"d",
		"",
		"Module description",
	)

	_ = moduleUploadCmd.MarkFlagRequired("title")
	_ = moduleUploadCmd.MarkFlagRequired("repr")

	return moduleUploadCmd
}
