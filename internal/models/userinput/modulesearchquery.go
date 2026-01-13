package userinput

type ModuleSearchQuery struct {
	Tags        []string
	Name        []string
	Description []string
}

func NewModuleSearchQuery() *ModuleSearchQuery {
	return &ModuleSearchQuery{}
}

func (inst ModuleSearchQuery) IsEmpty() bool {
	return len(inst.Description) == 0 && len(inst.Name) == 0 && len(inst.Tags) == 0
}
