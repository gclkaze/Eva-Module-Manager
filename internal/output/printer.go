package output

import "emm/internal/models"

type Printer interface {
	Info(msg string)
	Error(err error)
	PrintModules(mod []models.Module)
}
