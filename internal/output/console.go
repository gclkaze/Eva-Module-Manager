package output

import (
	"fmt"
	"os"
)

type ConsolePrinter struct{}

func NewConsolePrinter() *ConsolePrinter {
	return &ConsolePrinter{}
}

func (p *ConsolePrinter) Info(msg string) {
	fmt.Println(msg)
}

func (p *ConsolePrinter) Error(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}
