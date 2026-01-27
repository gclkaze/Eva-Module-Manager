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
		Args:  installArgs(&path),
		RunE:  installRunE(application, &path),
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "Path to eva.json, or a directory containing eva.json (default: ./eva.json)")
	_ = cmd.MarkFlagFilename("path", "json")

	return cmd
}

func installArgs(path *string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		pathSet := cmd.Flags().Changed("path")
		if len(args) > 1 {
			application.SetOnError()
			application.GetPrinter().Error(fmt.Errorf("invalid usage: expected 0 or 1 argument ([module@version])"))
			return nil
		}

		if pathSet && strings.TrimSpace(*path) == "" {
			application.SetOnError()
			application.GetPrinter().Error(fmt.Errorf("--path cannot be empty"))
			return nil
		}

		if len(args) == 1 {
			res, err := utils.IsValidModuleOrModuleVersion(args[0])
			if !res {
				application.SetOnError()
				application.GetPrinter().Error(err)
				return nil
			}
			if _, _, err := utils.ParseModuleReleaseVersion(args[0]); err != nil {
				application.GetPrinter().Error(err)
				application.SetOnError()
				return nil
			}
		}

		return nil
	}
}

func installRunE(application *app.EMMApp, path *string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if application.IsOnError() {
			return nil
		}
		pathSet := cmd.Flags().Changed("path")

		token, err := application.GetCurrentUserToken()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		ctx := cmd.Context()
		var summary *models.InstallationSummary

		//if module is absent, we install from eva.json or path/eva.json
		//if module is there, we install the module to eva.json or path/eva.json
		if len(args) == 0 {
			if pathSet {
				summary, err = installFromPath(ctx, token, application, *path)
			} else {
				summary, err = installAllFromProject(ctx, token, application)
			}
		} else {
			summary, err = installSingleModule(ctx, token, application, args[0], path)
		}
		if err != nil {
			application.GetPrinter().Error(err)
			application.GetPrinter().Error(fmt.Errorf("installation operation failed"))
			return nil
		}

		application.GetPrinter().PrintSummary(summary)

		//application.GetPrinter().Info("installation operation completed successfully.")
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

func installSingleModule(ctx context.Context, token string, application *app.EMMApp, moduleAtVersion string, path *string) (*models.InstallationSummary, error) {
	module, version, err := utils.ParseModuleReleaseVersion(moduleAtVersion)
	if err != nil {
		application.GetPrinter().Error(err)
		return nil, err
	}
	summary, err := application.InstallModuleVersion(ctx, token, module, version, path)
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
