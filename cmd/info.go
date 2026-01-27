package cmd

import (
	"emm/internal/app"
	"emm/pkg/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func NewShowModuleInfoCommand(application *app.EMMApp) *cobra.Command {
	module := ""

	var showCmd = &cobra.Command{
		Use:   "info",
		Short: "Show module or module release information",
		Long:  "Show module information or module release information",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				application.SetOnError()
				return fmt.Errorf("provide the module name or module-name@version for module/release information")
			}
			module = args[0]
			isValid, err := utils.IsValidModuleOrModuleVersion(module)
			if !isValid {
				application.SetOnError()
				application.GetPrinter().Error(err)
				return nil
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if application.IsOnError() {
				return nil
			}
			err := application.GetModuleInfo(module)
			if err != nil {
				application.GetPrinter().Error(err)
			}
			return nil
		},
	}

	return showCmd

}
