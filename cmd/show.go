package cmd

import "github.com/spf13/cobra"

var modules []string
var releases []string

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show module or module release information",
	Long:  "Show module information or module release information",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	showCmd.Flags().StringSliceVarP(
		&modules,
		"module",
		"m",
		[]string{},
		"Modules to search for (space-separated or repeatable)",
	)

	showCmd.Flags().StringSliceVarP(
		&releases,
		"release",
		"r",
		[]string{},
		"Module release to search for (space-separated or repeatable)",
	)

	// Optional: enforce at least one tag
	//searchCmd.MarkFlagRequired("tags")
}
