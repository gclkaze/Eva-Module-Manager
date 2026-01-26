package cmd

import (
	"emm/internal/app"
	"emm/pkg/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func NewReleaseDownloadCommand(application *app.EMMApp) *cobra.Command {
	var saveLocation string = ""
	var releaseDownloadCmd = &cobra.Command{
		Use:   "download module@version",
		Args:  cobra.ExactArgs(1),
		Short: "Download an available Module Release",
		RunE: func(cmd *cobra.Command, args []string) error {
			module, version, err := utils.ParseModuleReleaseVersion(args[0])
			if version == "" {
				version = "latest"
			}
			if err != nil {
				application.GetPrinter().Error(err)
				return nil
			}
			if saveLocation == "" {
				saveLocation = application.GetDefaultFileStorageLocation()
			}
			if !handleLocation(saveLocation) {
				return nil
			}
			tk, err := application.GetCurrentUserToken()
			if err != nil {
				//Anyone can download our Modules! But only admins can download any release with any status
				tk = ""
			}
			ctx := cmd.Context()
			err = application.DownloadRelease(ctx, tk, module, version, saveLocation)
			if err != nil {
				application.GetPrinter().Error(err)
			}
			return nil

		},
	}

	releaseDownloadCmd.Flags().StringVarP(
		&saveLocation,
		"savelocation",
		"s",
		"",
		"Folder where the release will be saved (must already exist)",
	)
	return releaseDownloadCmd

}

func handleLocation(saveLocation string) bool {
	if !utils.FolderExists(saveLocation) {
		err := utils.CreateFolder(saveLocation)
		if err != nil {
			application.GetPrinter().Error(err)
			return false
		}
		application.GetPrinter().Info(fmt.Sprintf("Created folder: '%s'.", saveLocation))
		if !utils.FolderExists(saveLocation) {
			return false
		}
	}
	return true
}
