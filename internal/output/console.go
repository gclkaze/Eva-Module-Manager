package output

import (
	"emm/internal/models"
	"emm/pkg/utils"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
)

type ConsolePrinter struct{}

func NewConsolePrinter() *ConsolePrinter {
	return &ConsolePrinter{}
}

func (p *ConsolePrinter) Info(msg string) {
	fmt.Println(msg)
}

func (p *ConsolePrinter) Error(err error) {
	errorColor := color.New(color.FgHiRed, color.Bold).SprintFunc()
	theError := err.Error()

	if strings.Contains(theError, "because the target machine actively refused it") {
		theError = "Couldn't connect to the Module Repository Server..check your internet connection"
	}
	fmt.Fprintln(os.Stderr, errorColor("Error: "+theError))
}

func (p *ConsolePrinter) Success(msg string) {
	successColor := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	fmt.Fprintln(os.Stderr, successColor(msg))
}

func (p ConsolePrinter) PrintModules(mods []models.Module) {
	maxRepr := 0
	for _, m := range mods {
		if len(m.Repr) > maxRepr {
			maxRepr = len(m.Repr)
		}
	}
	name := color.New(color.FgGreen).SprintFunc()
	meta := color.New(color.FgHiBlack).SprintFunc()

	for _, m := range mods {
		tags := ""
		if len(m.Tags) > 0 {
			sort.Strings(m.Tags)
			tags = p.formatTags(m.Tags, 5)
		}

		fmt.Printf(
			"%-*s (%s) - %s [%s]%s\n",
			maxRepr,
			name(m.RepoName),
			m.Title,
			m.Description,
			meta(fmt.Sprintf("[%d]", m.Releases)),
			meta(fmt.Sprintf("{%s}", tags)),
		)
	}
}

func (p ConsolePrinter) PrintModuleInfo(m models.ModuleEnrichedInformation) {
	maxRepr := 0
	if len(m.Repr) > maxRepr {
		maxRepr = len(m.Repr)
	}

	name := color.New(color.FgGreen).SprintFunc()
	meta := color.New(color.FgHiBlack).SprintFunc()

	tags := ""
	if len(m.Tags) > 0 {
		sort.Strings(m.Tags)
		tags = p.formatTags(m.Tags, 5)
	}

	fmt.Printf(
		"%-*s (%s) - %s [%s]%s\n",
		maxRepr,
		name(m.RepoName),
		m.Title,
		m.Description,
		meta(fmt.Sprintf("[%d]", len(m.ReleaseInfo))),
		meta(fmt.Sprintf("{%s}", tags)),
	)

	if len(m.ReleaseInfo) == 0 {
		error := color.New(color.BgHiRed).SprintFunc()
		fmt.Printf("Couldn't find any available releases for %s", error(m.RepoName))
		return
	}

	for i := range m.ReleaseInfo {
		p.PrintReleaseInfo(m.RepoName, m.ReleaseInfo[i])
	}

}

func (p ConsolePrinter) formatTags(tags []string, max int) string {
	if len(tags) <= max {
		return strings.Join(tags, ",")
	}
	return strings.Join(tags[:max], ",") + ",…"
}

func (p ConsolePrinter) PrintReleaseInfo(moduleRepr string, r models.Release) {
	version := color.New(color.FgCyan, color.Bold).SprintFunc()
	meta := color.New(color.FgHiBlack).SprintFunc()
	hint := color.New(color.FgYellow).SprintFunc()

	releasedAt := "N/A"
	if r.ReleasedAt != nil {
		releasedAt = r.ReleasedAt.Format("2006-01-02")
	}

	keywords := ""
	if len(r.Keywords) > 0 {
		labels := make([]string, len(r.Keywords))
		for i, k := range r.Keywords {
			labels[i] = k.Label
		}
		keywords = meta(fmt.Sprintf(" {%s}", strings.Join(labels, ",")))
	}

	fmt.Printf(
		"  └─ %s %s %s%s\n",
		version(r.Version),
		meta(releasedAt),
		meta(utils.HumanSize(r.DiskSize)),
		keywords,
	)

	if r.Description != "" {
		fmt.Printf("     %s\n", meta(r.Description))
	}

	fmt.Printf(
		"     %s %s\n",
		meta("→ install:"),
		hint(fmt.Sprintf(
			"emm install %s@%s",
			moduleRepr,
			r.Version,
		)),
	)
}
