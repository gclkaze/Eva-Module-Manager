package cmd

import (
	"emm/pkg/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var saveLocation string = ""
var releaseDownloadCmd = &cobra.Command{
	Use:   "download module@version",
	Args:  cobra.ExactArgs(1),
	Short: "Download an available Module Release",
	RunE: func(cmd *cobra.Command, args []string) error {
		module, version, err := utils.ParseModuleReleaseVersion(args[0])
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		if saveLocation == "" {
			saveLocation = application.GetDefaultFileStorageLocation()
		}
		if !checkLocation(saveLocation) {
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

func checkLocation(saveLocation string) bool {
	if saveLocation != "" {
		info, err := os.Stat(saveLocation)
		if err != nil {
			application.GetPrinter().Error(
				fmt.Errorf("invalid --savelocation: %w", err),
			)
			return false
		}
		if !info.IsDir() {
			application.GetPrinter().Error(
				fmt.Errorf("--savelocation must be a directory"),
			)
			return false
		}
	}
	return true
}

func init() {
	releaseCmd.AddCommand(releaseDownloadCmd)

	releaseDownloadCmd.Flags().StringVarP(
		&saveLocation,
		"savelocation",
		"s",
		"",
		"Folder where the release will be saved (must already exist)",
	)

}
