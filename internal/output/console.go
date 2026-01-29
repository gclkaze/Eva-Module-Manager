package output

import (
	"emm/internal/models"
	"emm/internal/models/dto"
	"emm/internal/models/eva"
	"emm/pkg/utils"
	"fmt"
	"os"
	"path/filepath"
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

type moduleRow struct {
	Info        eva.EvaModuleInfo
	AbsPath     string
	Referenced  bool
	OnDisk      bool
	DisplayPath string
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

func (p *ConsolePrinter) PrintSummary(s *models.InstallationSummary) {
	if s == nil {
		return
	}

	// Styles
	title := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	label := color.New(color.FgHiBlack).SprintFunc()

	ok := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	warn := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	fail := color.New(color.FgHiRed, color.Bold).SprintFunc()

	// Decide overall status
	statusText := "Completed"
	statusColor := ok
	if s.Failed > 0 {
		statusText = "Completed with errors"
		statusColor = fail
	} else if s.Skipped > 0 {
		statusText = "Completed with warnings"
		statusColor = warn
	}

	// Normalize processed (optional)
	processed := s.ProcessedCounter
	if processed <= 0 {
		processed = s.Success + s.Skipped + s.Failed
	}

	// Small, readable block with aligned values
	fmt.Println()
	fmt.Println(title("Installation Summary"), label("—"), statusColor(statusText))
	fmt.Println(strings.Repeat("─", 48))

	fmt.Printf("%-12s %s\n", label("Total:"), fmt.Sprintf("%d", s.Total))
	fmt.Printf("%-12s %s\n", label("Processed:"), fmt.Sprintf("%d", processed))

	fmt.Printf("%-12s %s\n", ok("Success:"), fmt.Sprintf("%d", s.Success))

	// Only show non-zero categories to reduce noise
	if s.Skipped > 0 {
		fmt.Printf("%-12s %s\n", warn("Skipped:"), fmt.Sprintf("%d", s.Skipped))
	}
	if s.Failed > 0 {
		fmt.Printf("%-12s %s\n", fail("Failed:"), fmt.Sprintf("%d", s.Failed))
	}

	fmt.Println(strings.Repeat("─", 48))

	// A short actionable hint line
	if !p.onVerboseMode {
		switch {
		case s.Failed > 0:
			fmt.Println(fail("Some modules failed."), label("Run with"), title("--verbose"), label("for details."))
		case s.Skipped > 0:
			fmt.Println(warn("Some modules were skipped."), label("Run with"), title("--verbose"), label("to see why."))
		default:
			fmt.Println(ok("All modules installed successfully."))
		}
	} else {
		switch {
		case s.Failed > 0:
			fmt.Println(fail("Some modules failed."))
		case s.Skipped > 0:
			fmt.Println(warn("Some modules were skipped."))
		default:
			fmt.Println(ok("All modules installed successfully."))
		}
	}
}

func (p ConsolePrinter) PrintDevelopers(
	devs []dto.DeveloperDTO,
	currentUserEmail string,
) {
	if len(devs) == 0 {
		color.New(color.FgHiBlack).Println("No developers found.")
		return
	}

	// Colors
	hdr := color.New(color.FgHiBlack, color.Bold).SprintFunc()
	nameC := color.New(color.FgGreen).SprintFunc()
	meta := color.New(color.FgHiBlack).SprintFunc()
	roleC := color.New(color.FgCyan).SprintFunc()
	bannedC := color.New(color.FgHiRed, color.Bold).SprintFunc()
	activeC := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	selfC := color.New(color.FgCyan, color.Bold).SprintFunc()

	const (
		wMark   = 1
		wID     = 6
		wBanned = 7
		wRole   = 10
		wEmail  = 34
		wHandle = 18
		wName   = 22
	)

	// Header
	fmt.Printf(
		"%s  %s  %s  %s  %s  %s  %s\n",
		hdr(""),
		hdr(fixedWidthTrunc("ID", wID)),
		hdr(fixedWidthTrunc("BANNED", wBanned)),
		hdr(fixedWidthTrunc("ROLE", wRole)),
		hdr(fixedWidthTrunc("EMAIL", wEmail)),
		hdr(fixedWidthTrunc("HANDLE", wHandle)),
		hdr("NAME"),
	)

	fixedColsWidth := wMark + wID + wBanned + wRole + wEmail + wHandle + wName
	gaps := 2 * 6
	fmt.Println(strings.Repeat("-", fixedColsWidth+gaps))

	currentUserEmail = strings.ToLower(strings.TrimSpace(currentUserEmail))

	for _, d := range devs {
		isSelf := strings.ToLower(strings.TrimSpace(d.Email)) == currentUserEmail

		mark := " "
		rowColor := func(s string) string { return s }
		if isSelf {
			mark = "*"
			rowColor = func(s string) string {
				return selfC(s)
			}
		}

		idStr := fixedWidthTrunc(fmt.Sprintf("%d", d.UserID), wID)

		statusStr := "ACTIVE"
		statusFn := activeC
		if d.IsBanned {
			statusStr = "BANNED"
			statusFn = bannedC
		}
		bannedStr := fixedWidthTrunc(statusStr, wBanned)

		role := strings.TrimSpace(d.UserRole)
		if role == "" {
			role = "-"
		}
		roleStr := fixedWidthTrunc(truncate(role, wRole), wRole)

		email := strings.TrimSpace(d.Email)
		if email == "" {
			email = "-"
		}
		emailStr := fixedWidthTrunc(truncate(email, wEmail), wEmail)

		handle := strings.TrimSpace(d.Handle)
		if handle == "" {
			handle = "-"
		}
		handleStr := fixedWidthTrunc(truncate(handle, wHandle), wHandle)

		name := strings.TrimSpace(d.FirstName + " " + d.LastName)
		if name == "" {
			name = handle
		}
		nameStr := fixedWidthTrunc(truncate(name, wName), wName)

		fmt.Printf(
			"%s  %s  %s  %s  %s  %s  %s\n",
			rowColor(mark),
			meta(idStr),
			statusFn(bannedStr),
			roleC(roleStr),
			rowColor(emailStr),
			meta(handleStr),
			nameC(nameStr),
		)
	}
}

func (p *ConsolePrinter) PrintUninstallSummary(s *models.PurgeSummary) {
	if s == nil {
		return
	}

	// Styles
	title := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	label := color.New(color.FgHiBlack).SprintFunc()

	ok := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	warn := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	fail := color.New(color.FgHiRed, color.Bold).SprintFunc()

	// Decide overall status
	statusText := "Completed"
	statusColor := ok
	if s.Failed > 0 {
		statusText = "Completed with errors"
		statusColor = fail
	} else if s.Skipped > 0 {
		statusText = "Completed with warnings"
		statusColor = warn
	}

	// Normalize processed (optional)
	processed := s.ProcessedCounter
	if processed <= 0 {
		processed = s.Success + s.Skipped + s.Failed
	}

	// Small, readable block with aligned values
	fmt.Println()
	fmt.Println(title("Uninstallation Summary"), label("—"), statusColor(statusText))
	fmt.Println(strings.Repeat("─", 48))

	fmt.Printf("%-12s %s\n", label("Total:"), fmt.Sprintf("%d", s.Total))
	fmt.Printf("%-12s %s\n", label("Processed:"), fmt.Sprintf("%d", processed))

	fmt.Printf("%-12s %s\n", ok("Success:"), fmt.Sprintf("%d", s.Success))

	// Only show non-zero categories to reduce noise
	if s.Skipped > 0 {
		fmt.Printf("%-12s %s\n", warn("Skipped:"), fmt.Sprintf("%d", s.Skipped))
	}
	if s.Failed > 0 {
		fmt.Printf("%-12s %s\n", fail("Failed:"), fmt.Sprintf("%d", s.Failed))
	}

	fmt.Println(strings.Repeat("─", 48))

	// A short actionable hint line
	if !p.onVerboseMode {
		switch {
		case s.Failed > 0:
			fmt.Println(fail("Some modules failed."), label("Run with"), title("--verbose"), label("for details."))
		case s.Skipped > 0:
			fmt.Println(warn("Some modules were skipped."), label("Run with"), title("--verbose"), label("to see why."))
		default:
			fmt.Println(ok("All modules installed successfully."))
		}
	} else {
		switch {
		case s.Failed > 0:
			fmt.Println(fail("Some modules failed."))
		case s.Skipped > 0:
			fmt.Println(warn("Some modules were skipped."))
		default:
			fmt.Println(ok("All modules installed successfully."))
		}
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

	if strings.Contains(theError, "An existing connection was forcibly closed by the remote host.") {
		theError = "The Module Repository Server went away..check your internet connection"
	}
	fmt.Fprintln(os.Stderr, errorColor("Error: "+theError))
}

func (p *ConsolePrinter) Success(msg string) {
	successColor := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	fmt.Fprintln(os.Stderr, successColor(msg))
}

func (p ConsolePrinter) PrintInstallationSummary() {

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

func (inst ConsolePrinter) PrintEvaModulesWithShowAll(p *eva.EvaProject, projectRoot, projectFileAbs string, showAll bool) {
	title := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	nameOK := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	nameMissing := color.New(color.FgHiRed, color.Bold).SprintFunc()
	verC := color.New(color.FgHiCyan).SprintFunc()
	pathC := color.New(color.FgHiBlack).SprintFunc()
	warnC := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	orphanC := color.New(color.FgHiYellow).SprintFunc() // on disk but not referenced

	// Always show where we loaded from
	if projectFileAbs != "" {
		fmt.Fprintf(os.Stdout, "%s %s\n", title("Project file:"), pathC(projectFileAbs))
	}
	fmt.Fprintln(os.Stdout)

	if p == nil {
		fmt.Fprintln(os.Stdout, warnC("No project loaded."))
		return
	}

	// 1) Collect referenced modules from eva.json
	referenced := map[string]eva.EvaModuleInfo{}
	for k, m := range p.Modules {
		referenced[k] = m
	}

	// 2) Scan filesystem modules if showAll
	onDisk := map[string]eva.EvaModuleInfo{}
	// expected structure: <projectRoot>/<ModulesFolder>/<moduleName>/<version>/
	base := filepath.Join(projectRoot, p.ModulesFolder)
	entries, err := os.ReadDir(base)
	if err == nil { // if folder missing, just treat as empty
		for _, modDir := range entries {
			if !modDir.IsDir() {
				continue
			}
			modName := modDir.Name()
			verBase := filepath.Join(base, modName)

			verEntries, err := os.ReadDir(verBase)
			if err != nil {
				continue
			}
			for _, verDir := range verEntries {
				if !verDir.IsDir() {
					continue
				}
				ver := verDir.Name()
				key := modName + "@" + ver
				relInstall := filepath.Join(p.ModulesFolder, modName, ver)

				onDisk[key] = eva.EvaModuleInfo{
					ModuleName:         modName,
					Version:            ver,
					InstallationFolder: relInstall,
				}
			}
		}
	}

	// 3) Build union rows
	unionKeys := make(map[string]struct{})
	for k := range referenced {
		unionKeys[k] = struct{}{}
	}
	if showAll {
		for k := range onDisk {
			unionKeys[k] = struct{}{}
		}
	}

	// If not showAll and no referenced modules
	if !showAll && len(unionKeys) == 0 {
		fmt.Fprintln(os.Stdout, warnC("No modules referenced in eva.json."))
		return
	}

	rows := make([]moduleRow, 0, len(unionKeys))
	for k := range unionKeys {
		var info eva.EvaModuleInfo
		ref := false
		disk := false

		if m, ok := referenced[k]; ok {
			info = m
			ref = true
		}
		if m, ok := onDisk[k]; ok {
			// prefer referenced info if present; otherwise use disk-derived info
			if !ref {
				info = m
			}
			disk = true
		}

		abs := info.InstallationFolder
		if projectRoot != "" {
			abs = filepath.Join(projectRoot, info.InstallationFolder)
		}

		rows = append(rows, moduleRow{
			Info:       info,
			AbsPath:    abs,
			Referenced: ref,
			OnDisk:     disk,
		})
	}

	// 4) Sort rows: name, version
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Info.ModuleName == rows[j].Info.ModuleName {
			return rows[i].Info.Version < rows[j].Info.Version
		}
		return rows[i].Info.ModuleName < rows[j].Info.ModuleName
	})

	// 5) Column widths
	maxName, maxVer := len("NAME"), len("VERSION")
	for _, r := range rows {
		if len(r.Info.ModuleName) > maxName {
			maxName = len(r.Info.ModuleName)
		}
		if len(r.Info.Version) > maxVer {
			maxVer = len(r.Info.Version)
		}
	}

	// 6) Header
	header := "Referenced modules"
	if showAll {
		header = "Modules (referenced + filesystem)"
	}
	fmt.Fprintln(os.Stdout, title(header))
	fmt.Fprintln(os.Stdout, strings.Repeat("─", maxName+maxVer+6+50))
	fmt.Fprintf(os.Stdout, "%-*s  %-*s  %s\n", maxName, "NAME", maxVer, "VERSION", "PATH")
	fmt.Fprintln(os.Stdout, strings.Repeat("─", maxName+maxVer+6+50))

	// 7) Rows with coloring rules
	var missingCount, orphanCount int

	for _, r := range rows {
		name := r.Info.ModuleName

		// Referenced but missing on disk => RED
		if r.Referenced && !r.OnDisk {
			name = nameMissing(name)
			missingCount++
		} else if !r.Referenced && r.OnDisk {
			// On disk but not referenced (only possible in showAll) => yellow (optional)
			name = orphanC(name)
			orphanCount++
		} else {
			// Normal referenced + present
			name = nameOK(name)
		}

		fmt.Fprintf(
			os.Stdout,
			"%-*s  %-*s  %s\n",
			maxName, name,
			maxVer, verC(r.Info.Version),
			pathC(r.AbsPath),
		)
	}

	// 8) Footer summary
	fmt.Fprintln(os.Stdout, strings.Repeat("─", maxName+maxVer+6+50))
	fmt.Fprintf(os.Stdout, "%s %d\n", title("Total:"), len(rows))

	if missingCount > 0 {
		fmt.Fprintf(os.Stdout, "%s %d\n", nameMissing("Missing on disk:"), missingCount)
	}
	if showAll && orphanCount > 0 {
		fmt.Fprintf(os.Stdout, "%s %d\n", warnC("On disk but not referenced:"), orphanCount)
	}
}
