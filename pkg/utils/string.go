package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rogpeppe/go-internal/semver"
)

var ErrInvalidModuleReleaseVersion = errors.New("invalid module-release version")

func HumanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf(
		"%.1f %cB",
		float64(bytes)/float64(div),
		"KMGTPE"[exp],
	)
}

func ParseModuleReleaseVersion(s string) (string, string, error) {
	parts := strings.Split(s, "@")

	switch len(parts) {
	case 1:
		// "X"
		if parts[0] == "" {
			return "", "", ErrInvalidModuleReleaseVersion
		}
		return parts[0], "", nil

	case 2:
		// "X@Y"
		module := parts[0]
		release := parts[1]

		if module == "" || release == "" {
			return "", "", ErrInvalidModuleReleaseVersion
		}

		if release == "latest" {
			return module, "latest", nil
		}

		// Normalize semver: must start with 'v'
		if !strings.HasPrefix(release, "v") {
			release = "v" + release
		}

		if !semver.IsValid(release) {
			return "", "", ErrInvalidModuleReleaseVersion
		}

		return module, release, nil

	default:
		// "X@Y@Z..."
		return "", "", ErrInvalidModuleReleaseVersion
	}
}
