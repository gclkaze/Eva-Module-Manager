package cmd

import (
	"emm/internal/app"
	"emm/internal/models"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewUninstallCommand(application *app.EMMApp) *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "uninstall [module@version]",
		Short: "Uninstalls an EVA module@version from eva.json.",
		Long:  "Uninstalls an EVA module@version from the local eva.json or an eva.json, that its location is provided through --path).",
		Args:  UninstallArgs(application, &path),
		RunE:  UninstallRunE(application, &path),
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "Path to eva.json, or a directory containing eva.json (default: ./eva.json)")
	_ = cmd.MarkFlagFilename("path", "json")

	return cmd
}

func UninstallArgs(application *app.EMMApp, path *string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		pathSet := cmd.Flags().Changed("path")
		if len(args) != 1 {
			return fmt.Errorf("invalid usage: expected 1 argument ([module or module@version])")
		}

		if pathSet && strings.TrimSpace(*path) == "" {
			return fmt.Errorf("--path cannot be empty")
		}
		return nil
	}
}

func UninstallRunE(application *app.EMMApp, path *string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		//pathSet := cmd.Flags().Changed("path")

		token, err := application.GetCurrentUserToken()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		ctx := cmd.Context()
		var summary *models.PurgeSummary

		summary, err = application.UninstallModule(ctx, token, args[0], path)

		if err != nil {
			application.GetPrinter().Error(err)
			application.GetPrinter().Error(fmt.Errorf("uninstallation operation failed"))
			return nil
		}

		application.GetPrinter().PrintUninstallSummary(summary)

		application.GetPrinter().Info("uninstallation operation completed successfully.")
		return nil
	}
}
