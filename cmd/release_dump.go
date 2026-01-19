package cmd

import (
	"emm/internal/models/userinput"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	releaseFilter    = userinput.NewReleaseFilterParams()
	createdAfterStr  string
	releasedAfterStr string
)

var releaseDumpCmd = &cobra.Command{
	Use:   "dump [params...]",
	Short: "Dump Module Release information, filter-powered",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return bindReleaseFilterTimes()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := application.GetCurrentUserToken()
		if err != nil {
			application.GetPrinter().Error(err)
			return nil
		}
		err = application.ReleaseDump(token, releaseFilter)
		if err != nil {
			application.GetPrinter().Error(err)
		}
		return nil
	},
}

func init() {
	releaseCmd.AddCommand(releaseDumpCmd)
	releaseFilter = userinput.NewReleaseFilterParams()
	f := releaseDumpCmd.Flags()

	// []string fields
	f.StringSliceVarP(&releaseFilter.Status, "status", "s", nil, "Filter by status")
	f.StringSliceVarP(&releaseFilter.Versions, "versions", "v", nil, "Filter by versions")
	f.StringSliceVarP(&releaseFilter.Tags, "tags", "t", nil, "Filter by tags")
	f.StringSliceVarP(&releaseFilter.ModuleName, "module", "m", nil, "Filter by module name(s)")
	f.StringSliceVarP(&releaseFilter.RepoName, "repo", "r", nil, "Filter by repo name(s)")
	f.StringSliceVarP(&releaseFilter.Description, "description", "d", nil, "Filter by description term(s)")
	f.StringSliceVarP(&releaseFilter.Creator, "creator", "c", nil, "Filter by creator(s)")
	f.StringSliceVarP(&releaseFilter.CreatorEmail, "creator-email", "e", nil, "Filter by creator email(s)")

	// time fields (string → parsed)
	f.StringVarP(&createdAfterStr, "created-after", "a", "", "Filter releases created after RFC3339 time")
	f.StringVarP(&releasedAfterStr, "released-after", "R", "", "Filter releases released after RFC3339 time")

}
func bindReleaseFilterTimes() error {
	if createdAfterStr != "" {
		t, err := time.Parse(time.RFC3339, createdAfterStr)
		if err != nil {
			return fmt.Errorf("--created-after must be RFC3339: %w", err)
		}
		releaseFilter.CreatedAfter = &t
	}

	if releasedAfterStr != "" {
		t, err := time.Parse(time.RFC3339, releasedAfterStr)
		if err != nil {
			return fmt.Errorf("--released-after must be RFC3339: %w", err)
		}
		releaseFilter.ReleasedAfter = t
	}

	return nil
}
