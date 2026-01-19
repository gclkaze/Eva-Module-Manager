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

func (p *ConsolePrinter) Warn(msg string) {
	warningColor := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	fmt.Fprintln(os.Stderr, warningColor(msg))
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
func (p ConsolePrinter) PrintReleaseRows(mods []models.ModuleEnrichedDTO) {
	if len(mods) == 0 {
		color.New(color.FgHiBlack).Println("No releases found.")
		return
	}

	// Flatten: one row per module@version
	type row struct {
		ModuleRepr string
		RepoName   string
		Release    models.ReleaseDTO
	}

	rows := make([]row, 0, 64)
	for _, m := range mods {
		for _, r := range m.ReleaseInfo {
			rows = append(rows, row{
				ModuleRepr: m.Repr,
				RepoName:   m.RepoName,
				Release:    r,
			})
		}
	}

	if len(rows) == 0 {
		color.New(color.FgHiBlack).Println("No releases found.")
		return
	}

	// Colors
	hdr := color.New(color.FgHiBlack, color.Bold).SprintFunc()
	moduleC := color.New(color.FgGreen).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()
	meta := color.New(color.FgHiBlack).SprintFunc()

	// Column widths (ASCII-aligned)
	const (
		wModule  = 22
		wID      = 30
		wStatus  = 10
		wRelTime = 16
		wSize    = 10
		wTags    = 22
		wDesc    = 34
	)

	// Header
	fmt.Printf(
		"%s  %s  %s  %s  %s  %s  %s\n",
		padRight(hdr("MODULE"), wModule),
		padRight(hdr("REPO@VERSION"), wID),
		padRight(hdr("STATUS"), wStatus),
		padRight(hdr("RELEASED"), wRelTime),
		padRight(hdr("SIZE"), wSize),
		padRight(hdr("TAGS"), wTags),
		padRight(hdr("DESCRIPTION"), wDesc),
	)

	fmt.Println(strings.Repeat("-", wModule+wID+wStatus+wRelTime+wSize+wTags+wDesc+12))

	// Rows
	for _, x := range rows {
		r := x.Release

		copyID := fmt.Sprintf("%s@%s", x.RepoName /*normalizeVersion*/, (r.Version))

		released := "-"
		if r.ReleasedAt != nil {
			released = r.ReleasedAt.Format("2006-01-02 15:04")
		}

		// Keywords -> comma separated labels
		tagLabels := make([]string, 0, len(r.Keywords))
		for _, k := range r.Keywords {
			tagLabels = append(tagLabels, k.Label)
		}
		tags := truncate(strings.Join(tagLabels, ","), wTags)

		desc := truncate(r.Description, wDesc)

		statusFn := statusColorFunc(r.Status)

		fmt.Printf(
			"%s  %s  %s  %s  %s  %s  %s\n",
			padRight(moduleC(truncate(x.ModuleRepr, wModule)), wModule),
			padRight(white(truncate(copyID, wID)), wID), // white copy-paste string
			padRight(statusFn(truncate(r.Status, wStatus)), wStatus),
			padRight(meta(truncate(released, wRelTime)), wRelTime),
			padRight(meta(truncate(humanSize(r.DiskSize), wSize)), wSize),
			padRight(meta(tags), wTags),
			padRight(meta(desc), wDesc),
		)
	}
}

// NOTE: These are byte-based; good for ASCII module names.
// If you need proper unicode width alignment, use a runewidth implementation.
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(s) <= max {
		return s
	}
	if max <= 1 {
		return s[:max]
	}
	return s[:max-1] + "…"
}

func (p ConsolePrinter) PrintDetailedModuleReleaseInfo(mods []models.ModuleEnrichedDTO) {
	if len(mods) == 0 {
		color.New(color.FgHiBlack).Println("No releases found.")
		return
	}

	// colors
	moduleTitle := color.New(color.FgGreen, color.Bold).SprintFunc()
	moduleMeta := color.New(color.FgHiBlack).SprintFunc()
	releaseTitle := color.New(color.FgCyan, color.Bold).SprintFunc()
	label := color.New(color.FgHiBlack).SprintFunc()
	value := color.New(color.FgWhite).SprintFunc()
	tagColor := color.New(color.FgYellow).SprintFunc()
	indexColor := color.New(color.FgHiBlack).SprintFunc()

	// header summary
	fmt.Println(indexColor(fmt.Sprintf("Modules: %d", len(mods))))
	fmt.Println()

	for mi, m := range mods {
		// ---- Module header ----
		fmt.Println(indexColor(fmt.Sprintf("[%d/%d] ", mi+1, len(mods))) +
			moduleTitle(m.Repr) + " " +
			moduleMeta(fmt.Sprintf("(%s)", m.RepoName)),
		)

		fmt.Println(label("Title:      "), value(m.Title))
		fmt.Println(label("Description:"), value(m.Description))

		if len(m.Tags) > 0 {
			fmt.Println(label("Tags:       "), tagColor(strings.Join(m.Tags, ", ")))
		}

		relCount := len(m.ReleaseInfo)
		fmt.Println(label("Releases:   "), value(fmt.Sprintf("%d", relCount)))
		fmt.Println(strings.Repeat("-", 60))

		// ---- Releases ----
		if relCount == 0 {
			fmt.Println(indexColor("No releases for this module."))
			fmt.Println()
			continue
		}

		for ri, r := range m.ReleaseInfo {
			statusFn := statusColorFunc(r.Status)

			copyID := fmt.Sprintf("%s@%s", m.RepoName, normalizeVersion(r.Version))

			fmt.Println(
				indexColor(fmt.Sprintf("  (%d/%d) ", ri+1, relCount)) +
					releaseTitle("Version:") + " " +
					value(r.Version) + " " +
					color.New(color.FgWhite).SprintFunc()(copyID),
			)

			fmt.Println(label("    Status:     "), statusFn(r.Status))

			if r.ReleasedAt != nil {
				fmt.Println(label("    Released at:"), value(r.ReleasedAt.Format("2006-01-02 15:04")))
			}

			fmt.Println(label("    Size:       "), value(humanSize(r.DiskSize)))
			fmt.Println(label("    Description:"), value(r.Description))

			if len(r.Keywords) > 0 {
				kws := make([]string, 0, len(r.Keywords))
				for _, k := range r.Keywords {
					kws = append(kws, k.Label)
				}
				fmt.Println(label("    Keywords:   "), tagColor(strings.Join(kws, ", ")))
			}

			fmt.Println()
		}

		fmt.Println()
	}
}
func normalizeVersion(v string) string {
	if len(v) > 0 && (v[0] == 'v' || v[0] == 'V') {
		return v[1:]
	}
	return v
}
func humanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func statusColorFunc(status string) func(a ...interface{}) string {
	switch strings.ToLower(status) {
	case "draft":
		return color.New(color.FgHiBlack).SprintFunc()
	case "pending":
		return color.New(color.FgYellow).SprintFunc()
	case "accepted":
		return color.New(color.FgGreen, color.Bold).SprintFunc()
	case "rejected":
		return color.New(color.FgRed, color.Bold).SprintFunc()
	case "canceled", "cancelled":
		return color.New(color.FgMagenta).SprintFunc()
	default:
		return color.New(color.FgWhite).SprintFunc()
	}
}

func (p ConsolePrinter) PrintModules(mods []models.Module) {
	if len(mods) == 0 {
		p.Info("The list of Modules is empty")
		return
	}
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
