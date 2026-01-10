package output

import "emm/internal/models"

type Printer interface {
	Info(msg string)
	Error(err error)
	PrintModules(mod []models.Module)
	PrintModuleInfo(m models.ModuleEnrichedInformation)
	PrintReleaseInfo(moduleRepr string, r models.Release)
}
