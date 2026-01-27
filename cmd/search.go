package cmd

import (
	"emm/internal/app"
	"emm/internal/models/userinput"
	"emm/pkg/utils"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewSearchArtifactsCommand(application *app.EMMApp) *cobra.Command {
	var searchQuery *userinput.ModuleSearchQuery

	var searchCmd = &cobra.Command{
		Use:     "search",
		Aliases: []string{"s"},
		Short:   "Search artifacts",
		Long:    "Search artifacts using one or more tags",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && searchQuery.IsEmpty() {
				application.SetOnError()
				return fmt.Errorf("provide search terms either as arguments or via --tags/--name/--description")
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			if application.IsOnError() {
				return
			}

			query := ""
			if searchQuery.IsEmpty() {
				query = strings.Join(args, " ")

			} else {
				searchQuery.Normalize()
				err := searchQuery.IsValid()
				if err != nil {
					application.GetPrinter().Error(err)
					return
				}
			}
			if query != "" {
				toks, err := utils.ParseSearchPhrases(query)
				if err != nil {
					application.GetPrinter().Error(err)
					return
				}

				query = strings.Join(toks, ",")
				err = application.SearchByQuery(query)
				if err != nil {
					application.GetPrinter().Error(err)
					return
				}
				return
			}

			err := application.SearchBySearchQuery(searchQuery)
			if err != nil {
				application.GetPrinter().Error(err)
			}
		},
	}

	searchQuery = userinput.NewModuleSearchQuery()
	searchCmd.Flags().StringSliceVarP(
		&searchQuery.Tags,
		"tags",
		"t",
		[]string{},
		"Tags to search for (space-separated or repeatable)",
	)

	searchCmd.Flags().StringSliceVarP(
		&searchQuery.Name,
		"name",
		"n",
		[]string{},
		"Name to search for (space-separated or repeatable)",
	)

	searchCmd.Flags().StringSliceVarP(
		&searchQuery.Description,
		"description",
		"d",
		[]string{},
		"Description to search for (space-separated or repeatable)",
	)
	return searchCmd
}
