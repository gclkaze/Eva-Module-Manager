package utils

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

func ReadPassword(prompt string) (string, error) {
	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		return "", errors.New(
			"stdin is not a terminal; use --password-stdin",
		)
	}

	fmt.Fprint(os.Stderr, prompt)
	pwd, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)

	return string(pwd), err
}
