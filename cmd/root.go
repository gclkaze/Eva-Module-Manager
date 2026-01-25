package cmd

import (
	"emm/internal/app"
	"emm/internal/config"
	"emm/internal/output"
	"emm/internal/services"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "emm",
	Short: "EMM is a CLI tool",
	Long:  "EMM is a CLI tool that can search artifacts using tags",
}

var application *app.EMMApp
var verbose bool

func initApp() error {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}

	moduleSearchService := services.NewModuleSearchService()
	authService := services.NewAuthService()
	moduleService := services.NewModuleService(authService)
	releaseService := services.NewModuleReleaseService(authService)
	bookKeepingService, err := services.NewProjectBookkeepingService(cwd)
	installService := services.NewInstallService(cwd, bookKeepingService, releaseService)

	if err != nil {
		return nil
	}

	application = app.NewEMMApp(
		cwd,
		moduleSearchService,
		authService,
		moduleService,
		releaseService,
		bookKeepingService,
		installService,
		output.NewConsolePrinter(verbose),
	)

	err = application.Init()
	if err != nil {
		return err
	}

	authService.SetProperties(config.TheConfigReader.GetProperties())
	moduleService.SetProperties(config.TheConfigReader.GetProperties())
	releaseService.SetProperties(config.TheConfigReader.GetProperties())
	bookKeepingService.SetProperties(config.TheConfigReader.GetProperties())
	installService.SetProperties(config.TheConfigReader.GetProperties())

	return nil
}
func initCmd() {
	rootCmd.AddCommand(
		NewVerifyCommand(application),
		NewInstallCommand(application),
		NewUninstallCommand(application),
	)

	rootCmd.PersistentFlags().BoolVarP(
		application.GetPrinter().GetVerboseFlagPointer(),
		"verbose",
		"v",
		false,
		"Enable verbose/debug output",
	)
}

func Execute() {
	err := initApp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initCmd()

	if err = rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
