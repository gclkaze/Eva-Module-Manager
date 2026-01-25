package models

type PurgeSummary struct {
	Total            int
	Skipped          int
	Success          int
	Failed           int
	ProcessedCounter int
}

func NewPurgeSummary() *PurgeSummary {
	return &PurgeSummary{}
}
