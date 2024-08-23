package scaffold

import (
	"github.com/hashicorp/go-version"
	"github.com/rs/zerolog"
)

type Metadata struct {
	MinimumVersion string `yaml:"minimum_version"`
}

func (m Metadata) IsCompatible(l zerolog.Logger, current string) (bool, error) {
	if m.MinimumVersion == "" || m.MinimumVersion == "*" {
		return true, nil
	}

	if current == "dev" {
		l.Debug().Msg("current version is dev, skipping version check")
		return true, nil
	}

	currentVersion, err := version.NewVersion(current)
	if err != nil {
		l.Error().Err(err).Msg("failed to parse current version")
		return false, err
	}

	minimumVersion, err := version.NewVersion(m.MinimumVersion)
	if err != nil {
		l.Error().Err(err).Msg("failed to parse minimum version")
		return false, err
	}

	return currentVersion.GreaterThanOrEqual(minimumVersion), nil
}
