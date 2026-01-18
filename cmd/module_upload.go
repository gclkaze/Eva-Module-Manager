package cmd

import (
	"emm/internal/models/userinput"
	"fmt"

	"github.com/spf13/cobra"
)

var uploadedCreds *userinput.UploadParams

var moduleUploadCmd = &cobra.Command{
	Use:   "upload [params...]",
	Short: "Upload a module",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tk, err := application.GetCurrentUserToken()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}

		err = uploadedCreds.AllValid()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		err = application.UploadModule(tk, args, uploadedCreds)
		if err != nil {
			application.GetPrinter().Error(err)
		}
		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Upload called with params:")
		for i, arg := range args {
			fmt.Printf("  %d: %s\n", i+1, arg)
		}
	},
}

func init() {
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

	moduleCmd.AddCommand(moduleUploadCmd)

}
