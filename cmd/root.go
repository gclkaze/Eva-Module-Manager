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

func initApp() error {
	moduleSearchService := services.NewModuleSearchService()
	authService := services.NewAuthService()

	application = app.NewEMMApp(
		moduleSearchService,
		authService,
		output.NewConsolePrinter(),
	)

	err := application.Init()
	if err != nil {
		return err
	}
	authService.SetProperties(config.TheConfigReader.GetProperties())
	return nil
}

func Execute() {
	err := initApp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err = rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
