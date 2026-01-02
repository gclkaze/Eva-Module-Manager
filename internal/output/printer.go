package output

type Printer interface {
	Info(msg string)
	Error(err error)
}
