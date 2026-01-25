package output

import "emm/internal/models"

type Printer interface {
	Info(msg string)
	VerboseInfo(msg string)
	Warn(msg string)
	VerboseWarn(msg string)

	Error(err error)
	Success(msg string)
	PrintModules(mod []models.Module)
	PrintModuleInfo(m models.ModuleEnrichedInformation)
	PrintReleaseInfo(moduleRepr string, r models.Release)
	PrintSummary(*models.InstallationSummary)
	PrintUninstallSummary(*models.PurgeSummary)

	PrintDetailedModuleReleaseInfo(mods []models.ModuleEnrichedDTO)
	PrintReleaseRows(mods []models.ModuleEnrichedDTO)

	GetVerboseFlagPointer() *bool
}
