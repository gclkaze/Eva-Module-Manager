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

type ConsolePrinter struct {
	onVerboseMode bool
}

// Flatten: one row per module@version
type row struct {
	ModuleRepr string
	RepoName   string
	Release    models.ReleaseDTO
	Tags       []string
}

func NewConsolePrinter(verbose bool) *ConsolePrinter {
	return &ConsolePrinter{onVerboseMode: verbose}
}

func (p *ConsolePrinter) GetVerboseFlagPointer() *bool {
	return &p.onVerboseMode
}

func (p *ConsolePrinter) Info(msg string) {
	fmt.Println(msg)
}

func (p *ConsolePrinter) Warn(msg string) {
	warningColor := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	fmt.Fprintln(os.Stderr, warningColor(msg))
}

func (p *ConsolePrinter) VerboseInfo(msg string) {
	if p.onVerboseMode {
		fmt.Println(msg)
	}
}

func (p *ConsolePrinter) VerboseWarn(msg string) {
	if p.onVerboseMode {
		warningColor := color.New(color.FgHiYellow, color.Bold).SprintFunc()
		fmt.Fprintln(os.Stderr, warningColor(msg))
	}
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

	rows := make([]row, 0, 64)
	for _, m := range mods {
		for _, r := range m.ReleaseInfo {
			rows = append(rows, row{
				ModuleRepr: m.Repr,
				RepoName:   m.RepoName,
				Release:    r,
				Tags:       m.Tags,
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
		wStatus  = 8
		wRelTime = 16
		wSize    = 6
		wTags    = 22
		wDescMax = 80 // description max only (no fixed width)
	)

	wModule := clamp(geModuleNameMaxColumnLength(rows), 10, 60)
	wID := clamp(getRepoVersionMaxColumnLength(rows), 10, 80)

	fmt.Printf(
		"%s  %s  %s  %s  %s  %s  %s\n",
		hdr(fixedWidthTrunc("MODULE", wModule)),
		hdr(fixedWidthTrunc("REPO@VERSION", wID)),
		hdr(fixedWidthTrunc("STATUS", wStatus)),
		hdr(fixedWidthTrunc("RELEASED", wRelTime)),
		hdr(fixedWidthTrunc("SIZE", wSize)),
		hdr(fixedWidthTrunc("TAGS", wTags)),
		hdr("DESCRIPTION"), // no fixed width
	)
	fixedColsWidth := wModule + wID + wStatus + wRelTime + wSize + wTags
	gaps := 2 * 6 // two spaces between 7 columns -> 6 gaps before DESCRIPTION
	fmt.Println(strings.Repeat("-", fixedColsWidth+gaps))

	// Rows
	for _, x := range rows {
		r := x.Release

		copyID := fmt.Sprintf("%s@%s", x.RepoName, normalizeVersion(r.Version))

		released := "-"
		if r.ReleasedAt != nil {
			released = r.ReleasedAt.Format("2006-01-02 15:04")
		}

		// Keywords -> comma separated labels
		tagLabels := make([]string, 0, len(r.Keywords))
		for _, k := range r.Keywords {
			tagLabels = append(tagLabels, k.Label)
		}

		tagLabels = append(tagLabels, x.Tags...)

		tags := truncate(strings.Join(tagLabels, ","), wTags)

		statusFn := statusColorFunc(r.Status)

		moduleStr := fixedWidthTrunc(truncate(x.ModuleRepr, wModule), wModule)
		idStr := fixedWidthTrunc(truncate(copyID, wID), wID)
		statusStr := fixedWidthTrunc(truncate(r.Status, wStatus), wStatus)
		releasedStr := fixedWidthTrunc(truncate(released, wRelTime), wRelTime)
		sizeStr := fixedWidthTrunc(truncate(humanSize(r.DiskSize), wSize), wSize)
		tagsStr := fixedWidthTrunc(truncMax(tags, wTags), wTags) // fixed width
		descStr := truncMax(r.Description, wDescMax)             // max 50 only

		fmt.Printf(
			"%s  %s  %s  %s  %s  %s  %s\n",
			moduleC(moduleStr),
			white(idStr),
			statusFn(statusStr),
			meta(releasedStr),
			meta(sizeStr),
			meta(tagsStr),
			meta(descStr),
		)

	}
}

func getRepoVersionMaxColumnLength(rows []row) int {
	max := 0
	for i := range rows {
		s := fmt.Sprintf("%s@%s", rows[i].RepoName, normalizeVersion(rows[i].Release.Version))
		l := len(s)
		if l > max {
			max = l
		}
	}
	return max
}

func geModuleNameMaxColumnLength(rows []row) int {
	max := 0
	for i := range rows {
		l := len(rows[i].ModuleRepr)
		if l > max {
			max = l
		}
	}
	return max
}

func clamp(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func fixedWidthTrunc(s string, width int) string {
	if len(s) > width {
		s = s[:width]
	}
	if len(s) < width {
		s += strings.Repeat(" ", width-len(s))
	}
	return s
}

func truncMax(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	if max == 1 {
		return s[:1]
	}
	return s[:max-1] + "…"
}

func fixedWidth(s string, width int) string {
	if len(s) > width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
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
