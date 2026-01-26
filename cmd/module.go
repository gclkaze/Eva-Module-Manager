package cmd

import (
	"emm/internal/app"

	"github.com/spf13/cobra"
)

func NewModuleParentCommand(application *app.EMMApp) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module",
		Short: "Module-related commands",
	}

	cmd.AddCommand(
		NewModuleSuggestionCommand(application),
		NewModuleUpdateCommand(application),
		NewModuleUploadCommand(application),
		NewModuleUserGetCommand(application),
	)
	return cmd
}
