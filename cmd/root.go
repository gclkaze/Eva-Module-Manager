package cmd

import (
	"emm/internal/app"
	"emm/internal/backend/perms"
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
	superviseService := services.NewSupervisionService(authService)

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
		superviseService,
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
	superviseService.SetProperties(config.TheConfigReader.GetProperties())

	return nil
}
func initCmd() {
	rootCmd.AddCommand(
		/*		NewVerifyCommand(application),
				NewInstallCommand(application),
				NewUninstallCommand(application),

				NewShowModuleInfoCommand(application),
				NewSearchArtifactsCommand(application),

				NewLoginCommand(application),
				NewRegisterCommand(application),
				NewSwitchUserCommand(application),
				NewLogoutCommand(application),
				NewWhoamiCommand(application),

				NewUserParentCommand(application),
				NewModuleParentCommand(application),
				NewReleaseParentCommand(application),*/

		enableFunctionsBasedOnPermissions(application, application.GetCurrentUserPermissions())...,
	)

	//perms := application.GetCurrentUserPermissions()

	rootCmd.PersistentFlags().BoolVarP(
		application.GetPrinter().GetVerboseFlagPointer(),
		"verbose",
		"v",
		false,
		"Enable verbose/debug output",
	)
}

func enableFunctionsBasedOnPermissions(application *app.EMMApp, thePerms map[string]bool) []*cobra.Command {
	var cmds []*cobra.Command
	cmds = append(cmds, NewVerifyCommand(application))
	cmds = append(cmds, NewInstallCommand(application))
	cmds = append(cmds, NewUninstallCommand(application))

	cmds = append(cmds, NewShowModuleInfoCommand(application))
	cmds = append(cmds, NewSearchArtifactsCommand(application))

	cmds = append(cmds, NewLoginCommand(application))
	cmds = append(cmds, NewRegisterCommand(application))
	cmds = append(cmds, NewSwitchUserCommand(application))
	cmds = append(cmds, NewLogoutCommand(application))
	cmds = append(cmds, NewWhoamiCommand(application))

	contains := perms.ContainsMappedPerms([]perms.Permission{perms.BanUsers, perms.UnbanUsers}, thePerms)
	if contains {
		cmds = append(cmds, NewUserParentCommand(application))
	}

	cmds = append(cmds, NewModuleParentCommand(application))
	cmds = append(cmds, NewReleaseParentCommand(application, thePerms))
	return cmds
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
