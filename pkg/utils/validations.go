package utils

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/rogpeppe/go-internal/semver"
)

var domainRegex = regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var nameRegex = regexp.MustCompile(`^[\p{L}]+([\p{L}\s'-]*[\p{L}]+)?$`)
var handleRegex = regexp.MustCompile(
	`^[a-zA-Z][a-zA-Z0-9]*(?:[._-][a-zA-Z0-9]+)*$`,
)
var reprRegex = regexp.MustCompile(`^[A-Za-z0-9]+( [A-Za-z0-9]+)*$`)
var ModuleReprMin = 3
var ModuleReprMax = 50

var passwordLengthMin = 12
var passwordLengthMax = 128

var nameMin = 5
var nameMax = 50

var handleMin = 3
var handleMax = 30

var repoNameMin = ModuleReprMin
var repoNameMax = ModuleReprMax

func IsValidEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return false
	}

	domain := parts[1]
	return domainRegex.MatchString(domain)
}
func IsValidName(name string) bool {
	if len(name) < nameMin || len(name) > nameMax {
		return false
	}
	return nameRegex.MatchString(name)
}

var repoNameRegex = regexp.MustCompile(
	`^[A-Za-z0-9][A-Za-z0-9._-]{0,98}[A-Za-z0-9]$`,
)

func IsValidRepoName(name string) error {
	if len(name) < ModuleReprMin {
		return fmt.Errorf("a valid module name is between %d and %d character length, consist of digits, numbers,dash or/and underscore", repoNameMin, repoNameMax)
	}
	if len(name) > ModuleReprMax {
		return fmt.Errorf("a valid module name is between %d and %d character length, consist of digits, numbers,dash or/and underscore", repoNameMin, repoNameMax)
	}

	if !repoNameRegex.MatchString(name) {
		return fmt.Errorf("a valid module name is between %d and %d character length, consist of digits, numbers,dash or/and underscore", repoNameMin, repoNameMax)
	}
	return nil
}

func IsValidVersion(v string) bool {
	if v[0] == 'v' {
		return semver.IsValid(v)
	}
	return semver.IsValid("v" + v)
}

func IsValidModuleName(moduleName string) bool {
	if len(moduleName) < ModuleReprMin || len(moduleName) > ModuleReprMax {
		return false
	}
	return reprRegex.MatchString(moduleName)
}

func IsValidHandle(handle string) bool {
	if len(handle) < handleMin || len(handle) > handleMax {
		return false
	}
	return handleRegex.MatchString(handle)
}

func IsValidNameWithError(attr string, name string) error {
	if len(name) < nameMin {
		return fmt.Errorf("%s must be at least %d characters long", attr, nameMin)
	}

	if len(name) > nameMax {
		return fmt.Errorf("%s must be at least %d characters long", attr, nameMax)
	}
	res := nameRegex.MatchString(name)
	if !res {
		return fmt.Errorf("only letters are allowed in %s", attr)
	}
	return nil
}

func IsValidHandleWithError(attr, handle string) error {
	if len(handle) < handleMin {
		return fmt.Errorf("%s must be at least %d characters long", attr, handleMin)
	}

	if len(handle) > handleMax {
		return fmt.Errorf("%s must be at least %d characters long", attr, handleMax)
	}
	res := handleRegex.MatchString(handle)
	if !res {
		return fmt.Errorf("only letters and numbers are allowed in %s, where the first character is a letter for instance :'mike85','mike'", attr)
	}
	return nil
}

var commonPasswords = map[string]struct{}{
	"password":  {},
	"123456":    {},
	"123456789": {},
	"qwerty":    {},
	"letmein":   {},
	"admin":     {},
}

func ValidatePassword(password, email string) error {
	pw := strings.TrimSpace(password)

	length := utf8.RuneCountInString(pw)
	if length < passwordLengthMin {
		return fmt.Errorf("password must be at least %d characters long", passwordLengthMin)
	}
	if length > passwordLengthMax {
		return errors.New("password is too long")
	}

	lowerPw := strings.ToLower(pw)
	if _, ok := commonPasswords[lowerPw]; ok {
		return errors.New("password is too common")
	}

	if email != "" {
		user := strings.ToLower(strings.Split(email, "@")[0])
		if user != "" && strings.Contains(lowerPw, user) {
			return errors.New("password must not contain your email or username")
		}
	}

	return nil
}

func ModuleReleasesAreEquivalent(keyA string, keyB string) (bool, error) {
	moduleA, versionA, err := ParseModuleReleaseVersion(keyA)
	if err != nil {
		return false, err
	}
	moduleB, versionB, err := ParseModuleReleaseVersion(keyB)
	if err != nil {
		return false, err
	}

	moduleA = strings.TrimSpace(moduleA)
	moduleB = strings.TrimSpace(moduleB)

	if moduleA != moduleB {
		return false, nil
	}
	versionA = strings.TrimSpace(versionA)
	versionB = strings.TrimSpace(versionB)

	return VersionsAreEquivalent(versionA, versionB), nil
}

func ReleasesAreEquivalent(moduleA string, versionA string, moduleB string, versionB string) bool {
	moduleA = strings.TrimSpace(moduleA)
	moduleB = strings.TrimSpace(moduleB)

	if moduleA != moduleB {
		return false
	}
	versionA = strings.TrimSpace(versionA)
	versionB = strings.TrimSpace(versionB)

	return VersionsAreEquivalent(versionA, versionB)
}

func VersionsAreEquivalent(versionA string, versionB string) bool {
	versionA = strings.TrimSpace(versionA)
	versionB = strings.TrimSpace(versionB)

	if versionA == "" && versionB != "" {
		return false
	}
	if versionA != "" && versionB == "" {
		return false
	}
	if versionA == "" && versionB == "" {
		return true
	}
	if versionA[0] == 'v' {
		versionA = versionA[1:]
	}
	if versionB[0] == 'v' {
		versionB = versionB[1:]
	}
	return versionA == versionB
}
