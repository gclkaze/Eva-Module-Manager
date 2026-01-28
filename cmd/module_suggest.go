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

			if err := suggestionParams.AllValid(); err != nil {
				application.GetPrinter().Error(err)
				return
			}

			if err := application.SuggestModuleRelease(tk, suggestionParams); err != nil {
				application.GetPrinter().Error(err)
			}

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

	suggestionCmd.Flags().StringVarP(
		&suggestionParams.Description,
		"description",
		"d",
		"",
		"Release description (optional)",
	)

	suggestionCmd.Flags().StringVarP(
		&suggestionParams.TagsCSV,
		"tags",
		"k",
		"",
		"Comma-separated release tags (optional)",
	)

	suggestionCmd.Flags().BoolVar(
		&suggestionParams.InheritModuleTags,
		"inherit-module-tags",
		false,
		"Inherit tags from the module when creating the release",
	)
	_ = suggestionCmd.MarkFlagRequired("module-name")
	_ = suggestionCmd.MarkFlagRequired("version")

	return suggestionCmd
}
