package cmd

import (
	"emm/internal/app"

	"github.com/spf13/cobra"
)

func NewReleaseParentCommand(application *app.EMMApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release-related commands",
	}

	cmd.AddCommand(
		NewReleaseAcceptCommand(application),
		NewReleaseCancelCommand(application),
		NewReleaseDownloadCommand(application),
		NewReleaseDumpCommand(application),
		NewReleaseLowerCommand(application),
		NewReleaseRejectCommand(application),
	)
	return cmd
}
