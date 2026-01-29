package cmd

import (
	"emm/internal/app"
	"emm/pkg/utils"

	"github.com/spf13/cobra"
)

func NewListCommand(application *app.EMMApp) *cobra.Command {
	var showAll bool
	var path string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List EVA modules used and referenced by the project file.",
		Long:    "List EVA modules used and referenced by the project file.",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			path, err = utils.ValidateOptionalDirPath(path)
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}

			if err := application.ListModules(showAll, path); err != nil {
				application.GetPrinter().Error(err)
			}
		},
	}

	cmd.Flags().BoolVar(
		&showAll,
		"show-all",
		false,
		"List all modules and releases used by the project and referenced in the file system",
	)

	cmd.Flags().StringVarP(
		&path,
		"path",
		"p",
		"",
		"Optional path to search for modules (defaults to project root)",
	)

	return cmd
}
