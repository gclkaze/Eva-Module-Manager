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
var repoRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
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

var MinVersionLength = 5
var MaxVersionLength = 64

var MinEmailValidLength = 10
var MaxEmailValidLength = 64

const (
	MaxDescriptionLength = 2000
	ModuleReprMinLength  = 2
	ModuleReprMaxLength  = 64
)

var moduleReprRe = regexp.MustCompile(
	`^[A-Za-z0-9]+([ _-][A-Za-z0-9]+)*$`,
)

func IsValidRepr(input string) error {
	name := strings.TrimSpace(input)

	if name == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	if len(name) < ModuleReprMinLength {
		return fmt.Errorf(
			"module name too short (min %d characters)",
			ModuleReprMinLength,
		)
	}

	if len(name) > ModuleReprMaxLength {
		return fmt.Errorf(
			"module name too long (max %d characters)",
			ModuleReprMaxLength,
		)
	}

	// Reject control characters
	for _, r := range name {
		if r < 32 {
			return fmt.Errorf(
				"module name contains invalid control characters",
			)
		}
	}

	if !moduleReprRe.MatchString(name) {
		return fmt.Errorf(
			"module name may contain only letters, numbers, spaces, '-' and '_'",
		)
	}

	return nil
}

func IsValidDescription(desc string) error {
	// Allow empty description
	if desc == "" {
		return nil
	}

	// Trim outer whitespace for validation only
	s := strings.TrimSpace(desc)

	// Length check
	if len(s) > MaxDescriptionLength {
		return fmt.Errorf(
			"description too long (%d characters, max %d)",
			len(s),
			MaxDescriptionLength,
		)
	}

	// Control character check
	for i, r := range s {
		if r < 32 && r != '\n' && r != '\t' {
			return fmt.Errorf(
				"description contains invalid control character at position %d",
				i,
			)
		}
	}

	return nil
}

func IsValidModuleOrModuleVersion(s string) (bool, error) {
	module, version, err := ParseModuleReleaseVersion(s)
	if err != nil {
		return false, err
	}

	err = IsValidModuleNameWithError(module)
	if err != nil {
		return false, err
	}

	if version == "" || version == "latest" {
		return true, nil
	}
	if !IsValidVersion(version) {
		return false, fmt.Errorf("incorrect version: '%s'. Additionally its length needs to be between %d and %d", version, MinVersionLength, MaxVersionLength)
	}

	return true, nil
}

func IsValidSpecificModuleVersion(s string) (bool, error) {
	module, version, err := ParseModuleReleaseVersion(s)
	if err != nil {
		return false, err
	}

	if version == "" || version == "latest" {
		return false, fmt.Errorf("the release needs to be specific, not just '%s'", version)
	}

	err = IsValidModuleNameWithError(module)
	if err != nil {
		return false, err
	}

	if !IsValidVersion(version) {
		return false, fmt.Errorf("incorrect version: '%s'. Additionally its length needs to be between %d and %d", version, MinVersionLength, MaxVersionLength)
	}

	return true, nil
}

func IsValidEmail(email string) bool {
	if len(email) < MinEmailValidLength {
		return false
	}
	if len(email) > MaxEmailValidLength {
		return false
	}

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
	if len(v) > MaxVersionLength {
		return false
	}
	if len(v) < MinVersionLength {
		return false
	}

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

func IsValidModuleNameWithError(moduleName string) error {
	attr := "module name"
	if len(moduleName) < ModuleReprMin {
		return fmt.Errorf("%s must be at least %d characters long, it is %d", attr, ModuleReprMin, len(moduleName))
	}
	if len(moduleName) > ModuleReprMax {
		return fmt.Errorf("%s must be at least %d characters long, it is %d", attr, ModuleReprMax, len(moduleName))
	}

	res := reprRegex.MatchString(moduleName)
	if !res {
		res = repoRegex.MatchString(moduleName)
		if !res {
			return fmt.Errorf("only letters,digits and dashes are allowed in %s", attr)
		}
	}

	return nil
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

const (
	// Max number of terms per field
	MaxTags         = 10
	MaxNames        = 5
	MaxDescriptions = 5

	// Length limits
	TagMinLen = 1
	TagMaxLen = 32

	NameMinLen = 1
	NameMaxLen = 64

	DescMinLen = 1
	DescMaxLen = 128
)

var (
	// tags are identifiers
	tagRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	// name / description search terms (human text)
	textRe = regexp.MustCompile(`^[\p{L}\p{N} ._-]+$`)
)

func validateList(
	label string,
	values []string,
	maxItems int,
	minLen int,
	maxLen int,
	re *regexp.Regexp,
) error {
	// Empty slice is valid
	if len(values) == 0 {
		return nil
	}

	if len(values) > maxItems {
		return fmt.Errorf(
			"%s has too many values (max %d)",
			label,
			maxItems,
		)
	}

	for i, v := range values {
		if len(v) < minLen || len(v) > maxLen {
			return fmt.Errorf(
				"%s[%d] length must be between %d and %d",
				label,
				i,
				minLen,
				maxLen,
			)
		}

		// Reject control characters
		for _, r := range v {
			if r < 32 {
				return fmt.Errorf(
					"%s[%d] contains invalid control characters",
					label,
					i,
				)
			}
		}

		if !re.MatchString(v) {
			return fmt.Errorf(
				"%s[%d] contains invalid characters",
				label,
				i,
			)
		}
	}

	return nil
}

func NormalizeTags(in []string) []string {
	return normalizeTerms(in)
}

func NormalizeText(in []string) []string {
	return normalizeTerms(in)
}

func ValidateTags(tags []string) error {
	return validateList(
		"tags",
		tags,
		MaxTags,
		TagMinLen,
		TagMaxLen,
		tagRe,
	)
}

func ValidateNames(names []string) error {
	return validateList(
		"name",
		names,
		MaxNames,
		NameMinLen,
		NameMaxLen,
		textRe,
	)
}

func ValidateDescriptions(desc []string) error {
	return validateList(
		"description",
		desc,
		MaxDescriptions,
		DescMinLen,
		DescMaxLen,
		textRe,
	)
}

const (
	SearchRawMaxLen   = 512
	SearchMaxTokens   = 10
	SearchTokenMaxLen = 64
)

var searchTokenRe = regexp.MustCompile(`^[\p{L}\p{N}._-]+$`)

// ParseSearchTerms parses a comma/space mixed search string into normalized tokens.
// Empty input is allowed (returns empty slice).
func ParseSearchTerms(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	if len(raw) > SearchRawMaxLen {
		return nil, fmt.Errorf("search string too long (max %d characters)", SearchRawMaxLen)
	}

	// reject control chars early
	for _, r := range raw {
		if r < 32 && r != '\t' && r != '\n' {
			return nil, fmt.Errorf("search string contains invalid control characters")
		}
	}

	// Split by commas first (phrases)
	chunks := strings.Split(raw, ",")

	seen := make(map[string]struct{}, 16)
	out := make([]string, 0, 16)

	for _, c := range chunks {
		c = strings.TrimSpace(c)
		if c == "" {
			continue // allow extra commas but they contribute nothing
		}

		// Then split by whitespace to get consistent tokens
		parts := strings.Fields(c)
		for _, p := range parts {
			p = strings.ToLower(p)

			if len(p) > SearchTokenMaxLen {
				return nil, fmt.Errorf("search token too long (max %d characters): %q", SearchTokenMaxLen, p)
			}

			// token-level char whitelist (no spaces here anymore)
			if !searchTokenRe.MatchString(p) {
				return nil, fmt.Errorf("search token contains invalid characters: %q", p)
			}

			if _, ok := seen[p]; ok {
				continue
			}
			seen[p] = struct{}{}
			out = append(out, p)

			if len(out) > SearchMaxTokens {
				return nil, fmt.Errorf("too many search tokens (max %d)", SearchMaxTokens)
			}
		}
	}

	return out, nil
}

const (
	SearchPhrasesRawMaxLen = 512
	SearchPhrasesMaxItems  = 10
	SearchPhraseMinLen     = 1
	SearchPhraseMaxLen     = 128
)

var searchPhraseRe = regexp.MustCompile(`^[\p{L}\p{N}._-]+( [\p{L}\p{N}._-]+)*$`)

func ParseSearchPhrases(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	if len(raw) > SearchPhrasesRawMaxLen {
		return nil, fmt.Errorf("search string too long (max %d characters)", SearchPhrasesRawMaxLen)
	}

	// reject control characters early (keep it one-line friendly)
	for _, r := range raw {
		if r < 32 { // disallow \n, \t too for this field
			return nil, fmt.Errorf("search string contains invalid control characters")
		}
	}

	chunks := strings.Split(raw, ",")

	out := make([]string, 0, len(chunks))
	seen := make(map[string]struct{}, len(chunks))

	for _, c := range chunks {
		c = strings.TrimSpace(c)
		if c == "" {
			continue // allow extra commas, they just add no phrase
		}

		// normalize internal whitespace to single spaces
		phrase := strings.Join(strings.Fields(c), " ")
		phrase = strings.ToLower(phrase)

		if len(phrase) < SearchPhraseMinLen || len(phrase) > SearchPhraseMaxLen {
			return nil, fmt.Errorf(
				"search phrase length must be between %d and %d characters: %q",
				SearchPhraseMinLen,
				SearchPhraseMaxLen,
				phrase,
			)
		}

		if !searchPhraseRe.MatchString(phrase) {
			return nil, fmt.Errorf("search phrase contains invalid characters: %q", phrase)
		}

		if _, ok := seen[phrase]; ok {
			continue
		}
		seen[phrase] = struct{}{}
		out = append(out, phrase)

		if len(out) > SearchPhrasesMaxItems {
			return nil, fmt.Errorf("too many search phrases (max %d)", SearchPhrasesMaxItems)
		}
	}

	return out, nil
}

func ParseTagsCSV(input string) ([]string, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return nil, nil // whole field empty is allowed
	}

	parts := strings.Split(raw, ",")
	seen := make(map[string]struct{}, len(parts))
	out := make([]string, 0, len(parts))

	for i, p := range parts {
		tag := strings.ToLower(strings.TrimSpace(p))

		// 🚫 empty elements are NOT allowed
		if tag == "" {
			return nil, fmt.Errorf("tags[%d] is empty", i)
		}

		if len(tag) < TagMinLen || len(tag) > TagMaxLen {
			return nil, fmt.Errorf(
				"tags[%d] length must be between %d and %d characters",
				i,
				TagMinLen,
				TagMaxLen,
			)
		}

		if !tagRe.MatchString(tag) {
			return nil, fmt.Errorf("tags[%d] contains invalid characters", i)
		}

		if _, ok := seen[tag]; ok {
			continue // duplicates are fine, just ignored
		}

		seen[tag] = struct{}{}
		out = append(out, tag)

		if len(out) > MaxTags {
			return nil, fmt.Errorf("too many tags (max %d)", MaxTags)
		}
	}

	return out, nil
}
