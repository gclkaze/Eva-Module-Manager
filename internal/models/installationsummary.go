package models

type InstallationSummary struct {
	Total            int
	Skipped          int
	Success          int
	Failed           int
	ProcessedCounter int
}

func NewInstallationSummary() *InstallationSummary {
	return &InstallationSummary{}
}
