package cmd

import (
	"emm/internal/models/userinput"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

/*
var tags []string
var name []string
var description []string
*/
var searchQuery *userinput.ModuleSearchQuery

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search artifacts",
	Long:  "Search artifacts using one or more tags",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && searchQuery.IsEmpty() {
			application.SetOnError()
			return fmt.Errorf("provide search terms either as arguments or via --tags/--name/--description")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if application.IsOnError() {
			return nil
		}

		query := ""
		if searchQuery.IsEmpty() {
			query = strings.Join(args, " ")
			fmt.Println("Search query:", query)

		}
		if query != "" {
			err := application.SearchByQuery(query)
			return err
		}

		err := application.SearchBySearchQuery(searchQuery)
		return err
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchQuery = userinput.NewModuleSearchQuery()
	// --tags "super parser"
	// --tags super --tags parser
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

	// Optional: enforce at least one tag
	//searchCmd.MarkFlagRequired("tags")
}
