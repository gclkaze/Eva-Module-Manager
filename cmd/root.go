package cmd

import (
	"emm/internal/app"
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

func initApp() {
	/*	deployService := &services.DeployService{
		client: newAPIClient(),
	}*/

	moduleSearchService := services.NewModuleSearchService()

	application = app.NewEMMApp(
		moduleSearchService,
		//		Deploy: deployService,
		output.NewConsolePrinter(),
	)
}

func Execute() {
	initApp()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
