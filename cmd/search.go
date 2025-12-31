package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var tags []string

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search artifacts",
	Long:  "Search artifacts using one or more tags",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Allow: --tags "super parser"
		if len(tags) == 1 {
			tags = strings.Fields(tags[0])
		}

		fmt.Println("Searching with tags:")
		for _, tag := range tags {
			fmt.Println(" -", tag)
		}

		// Here you would call your service / repository layer
		// results := service.SearchByTags(tags)

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

	// Optional: enforce at least one tag
	searchCmd.MarkFlagRequired("tags")
}
