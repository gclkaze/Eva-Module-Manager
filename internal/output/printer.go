package output

import "emm/internal/models"

type Printer interface {
	Info(msg string)
	Warn(msg string)
	Error(err error)
	Success(msg string)
	PrintModules(mod []models.Module)
	PrintModuleInfo(m models.ModuleEnrichedInformation)
	PrintReleaseInfo(moduleRepr string, r models.Release)
}
