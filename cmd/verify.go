// cmd/verify.go
package cmd

import (
	"emm/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

type VerifyOptions struct {
	// Can be:
	//  - empty: default to ./eva.json
	//  - a directory: <dir>/eva.json
	//  - a file: must be eva.json
	Path string
}

func NewVerifyCommand(application *app.EMMApp) *cobra.Command {
	opts := &VerifyOptions{}

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify eva.json structure",
		Long:  "Validates eva.json schema, module keys, and path safety (supports floating module versions).",
		Run: func(cmd *cobra.Command, args []string) {
			// Verify must NOT create eva.json; it only checks if it exists and is valid.
			absPath, err := application.VerifyEvaProjectFile(opts.Path)
			if err != nil {
				application.GetPrinter().Error(err)
				return
			}
			application.GetPrinter().Info(fmt.Sprintf("✅ eva.json is OK (%s)", absPath))
		},
	}

	// Naming suggestion: --path (accepts either a directory or file).
	// Alternative: --file, if you only want file paths.
	cmd.Flags().StringVarP(&opts.Path, "path", "p", "", "Path to eva.json, or a directory containing eva.json (default: ./eva.json)")
	_ = cmd.MarkFlagFilename("path", "json")

	return cmd
}
