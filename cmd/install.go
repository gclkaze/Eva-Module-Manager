package cmd

import (
	"context"
	"fmt"
	"strings"

	"emm/internal/app"
	"emm/internal/models"
	"emm/pkg/utils"

	"github.com/spf13/cobra"
)

func NewInstallCommand(application *app.EMMApp) *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "install [module@version]",
		Short: "Downloads && installs EVA modules (all from eva.json, a single module@version, or from --path).",
		Long:  "Downloads && installs EVA modules (all from eva.json, a single module@version, or from --path).",
		Args:  installArgs(application, &path),
		RunE:  installRunE(application, &path),
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "Path to eva.json, or a directory containing eva.json (default: ./eva.json)")
	_ = cmd.MarkFlagFilename("path", "json")

	return cmd
}

func installArgs(application *app.EMMApp, path *string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		pathSet := cmd.Flags().Changed("path")

		if pathSet && len(args) > 0 {
			return fmt.Errorf("invalid usage: use either [module@version] or --path, not both")
		}

		if len(args) > 1 {
			return fmt.Errorf("invalid usage: expected 0 or 1 argument ([module@version])")
		}

		if pathSet && strings.TrimSpace(*path) == "" {
			return fmt.Errorf("--path cannot be empty")
		}

		if len(args) == 1 {
			if _, _, err := utils.ParseModuleReleaseVersion(args[0]); err != nil {
				return err
			}
		}

		return nil
	}
}

func installRunE(application *app.EMMApp, path *string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		pathSet := cmd.Flags().Changed("path")

		token, err := application.GetCurrentUserToken()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		ctx := cmd.Context()
		var summary *models.InstallationSummary
		switch {
		case pathSet:
			summary, err = installFromPath(ctx, token, application, *path)
		case len(args) == 1:
			summary, err = installSingleModule(ctx, token, application, args[0])
		default:
			summary, err = installAllFromProject(ctx, token, application)
		}

		if err != nil {
			application.GetPrinter().Error(err)
			application.GetPrinter().Error(fmt.Errorf("operation failed"))
			return nil
		}

		application.GetPrinter().PrintSummary(summary)

		application.GetPrinter().Info("install operation completed successfully.")
		return nil
	}
}

func installFromPath(ctx context.Context, token string, application *app.EMMApp, path string) (*models.InstallationSummary, error) {
	summary, err := application.InstallAllFromPath(ctx, token, path)
	if err != nil {
		application.GetPrinter().Error(err)
		return summary, err
	}
	return summary, nil
}

func installSingleModule(ctx context.Context, token string, application *app.EMMApp, moduleAtVersion string) (*models.InstallationSummary, error) {
	module, version, err := utils.ParseModuleReleaseVersion(moduleAtVersion)
	if err != nil {
		application.GetPrinter().Error(err)
		return nil, err
	}
	summary, err := application.InstallModuleVersion(ctx, token, module, version)
	if err != nil {
		application.GetPrinter().Error(err)
		return summary, err
	}
	return summary, nil
}

func installAllFromProject(ctx context.Context, token string, application *app.EMMApp) (*models.InstallationSummary, error) {
	summary, err := application.InstallAllFromProject(ctx, token)
	if err != nil {
		application.GetPrinter().Error(err)
		return summary, err
	}
	return summary, nil
}
