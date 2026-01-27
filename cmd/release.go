package cmd

import (
	"emm/internal/app"
	"emm/internal/backend/perms"

	"github.com/spf13/cobra"
)

func NewReleaseParentCommand(application *app.EMMApp, thePerms map[string]bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Release-related commands",
	}

	cmd.AddCommand(
		NewReleaseDownloadCommand(application),
	)

	if perms.ContainsMappedPerms([]perms.Permission{perms.UpdateReleases}, thePerms) {
		cmd.AddCommand(NewReleaseDumpCommand(application))
	}
	if perms.ContainsMappedPerms([]perms.Permission{perms.AcceptReleases}, thePerms) {
		cmd.AddCommand(NewReleaseAcceptCommand(application))
	}
	if perms.ContainsMappedPerms([]perms.Permission{perms.CancelReleases}, thePerms) {
		cmd.AddCommand(NewReleaseCancelCommand(application))
	}
	if perms.ContainsMappedPerms([]perms.Permission{perms.ChangeReleaseStatuses}, thePerms) {
		cmd.AddCommand(NewReleaseLowerCommand(application))
	}
	if perms.ContainsMappedPerms([]perms.Permission{perms.RejectReleases}, thePerms) {
		cmd.AddCommand(NewReleaseRejectCommand(application))
	}

	return cmd
}
