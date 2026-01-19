package cmd

import (
	"emm/internal/models/userinput"

	"github.com/spf13/cobra"
)

var suggestionParams *userinput.ModuleReleaseSuggestionParams

var suggestionCmd = &cobra.Command{
	Use:   "suggest [params...]",
	Short: "Suggest a module for release",
	RunE: func(cmd *cobra.Command, args []string) error {
		tk, err := application.GetCurrentUserToken()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}

		err = suggestionParams.AllValid()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		err = application.SuggestModuleRelease(tk, suggestionParams)
		if err != nil {
			application.GetPrinter().Error(err)
		}
		return nil

	},
}

func init() {
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
		"v",
		"",
		"Module release version (required)",
	)

	_ = suggestionCmd.MarkFlagRequired("module-name")
	_ = suggestionCmd.MarkFlagRequired("version")

	moduleCmd.AddCommand(suggestionCmd)

}
