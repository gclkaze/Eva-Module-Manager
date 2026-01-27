package cmd

import (
	"emm/internal/app"
	"emm/internal/models/userinput"

	"github.com/spf13/cobra"
)

func NewModuleSuggestionCommand(application *app.EMMApp) *cobra.Command {
	var suggestionParams *userinput.ModuleReleaseSuggestionParams

	var suggestionCmd = &cobra.Command{
		Use:   "suggest [params...]",
		Short: "Suggest a module for release",
		Run: func(cmd *cobra.Command, args []string) {
			tk, err := application.GetCurrentUserToken()
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}

			err = suggestionParams.AllValid()
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}
			err = application.SuggestModuleRelease(tk, suggestionParams)
			if err != nil {
				application.GetPrinter().Error(err)
			}
			return

		},
	}
	suggestionParams = userinput.NewModuleReleaseSuggestionParams()

	suggestionCmd.Flags().StringVarP(
		&suggestionParams.ModuleRepr,
		"module-name",
		"m",
		"",
		"Module name (required)",
	)

	suggestionCmd.Flags().StringVarP(
		&suggestionParams.Version,
		"version",
		"t",
		"",
		"Module release version (required)",
	)

	_ = suggestionCmd.MarkFlagRequired("module-name")
	_ = suggestionCmd.MarkFlagRequired("version")

	return suggestionCmd
}
