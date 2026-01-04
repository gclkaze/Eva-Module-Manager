package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var tags []string
var name []string
var description []string

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search artifacts",
	Long:  "Search artifacts using one or more tags",
	Args: func(cmd *cobra.Command, args []string) error {
		// Either flags OR positional args must be provided
		if len(args) == 0 && len(tags) == 0 && len(description) == 0 {
			return fmt.Errorf("provide search terms either as arguments or via --tags/--name/--description")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		// Allow: --tags "super parser"
		/*		if len(tags) == 1 {
					tags = strings.Fields(tags[0])
				}

				if len(name) == 1 {
					name = strings.Fields(name[0])
				}

				if len(description) == 1 {
					description = strings.Fields(description[0])
				}*/

		/*		fmt.Println("Searching with tags:")
				for _, tag := range tags {
					fmt.Println(" -", tag)
				}

				fmt.Println("Searching with name:")
				for _, n := range name {
					fmt.Println(" -", n)
				}

				fmt.Println("Searching with name:")
				for _, d := range description {
					fmt.Println(" -", d)
				}*/

		query := ""
		if len(description) == 0 && len(name) == 0 && len(tags) == 0 {
			query = strings.Join(args, " ")
			fmt.Println("Search query:", query)

		}
		if query != "" {
			err := application.SearchByQuery(query)
			return err
		} else {
			err := application.SearchByComponents(name, description, tags)
			return err

		}

		//is everything empty?
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// --tags "super parser"
	// --tags super --tags parser
	searchCmd.Flags().StringSliceVarP(
		&tags,
		"tags",
		"t",
		[]string{},
		"Tags to search for (space-separated or repeatable)",
	)

	searchCmd.Flags().StringSliceVarP(
		&name,
		"name",
		"n",
		[]string{},
		"Name to search for (space-separated or repeatable)",
	)

	searchCmd.Flags().StringSliceVarP(
		&description,
		"description",
		"d",
		[]string{},
		"Description to search for (space-separated or repeatable)",
	)

	// Optional: enforce at least one tag
	//searchCmd.MarkFlagRequired("tags")
}
