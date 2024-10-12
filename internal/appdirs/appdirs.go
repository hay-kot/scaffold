// Package appdirs provides a cross-platform way to find application directories for
// the configuration, cache, and any other application-specific directories.
//
// This package largely exists because of the legacy behavior of the scaffold CLI where
// we want to preserve user locations for configuration if it exists, but move forward
// to use the XDG Base Directory Specification for new installations.
package appdirs

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const (
	RCFilename = "scaffoldrc.yml"
)

// RCFilepath returns the path to the scaffold configuration file. This is done
// by:
//
//  1. Checking if the legacy configuration file exists in the user's home directory.
//  2. If it does, return the path to the legacy configuration file.
//  3. If it does not, return the path to the XDG configuration file.
func RCFilepath() string {
	legacyFilepath, exists := RCFilepathLegacy()
	if exists {
		return legacyFilepath
	}

	return RCFilepathXDG()
}

func RCFilepathLegacy() (path string, exists bool) {
	fp := homeDir(".scaffold", RCFilename)
	_, err := os.Stat(fp)
	if err != nil {
		return "", false
	}

	return fp, true
}

func RCFilepathXDG() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get home directory")
		}

		configDir = filepath.Join(home, ".config")
	}

	return filepath.Join(configDir, "scaffold", RCFilename)
}

// CacheDir returns the path to the scaffold cache directory. This is done by:
//
//  1. Checking if the XDG_DATA_HOME environment variable is set.
//  2. If it is, return the path to the cache directory.
//  3. If it is not, return the path to the cache directory in the user's home directory.
func CacheDir() string {
	legacyFilepath, exists := CacheDirLegacy()
	if exists {
		return legacyFilepath
	}

	return CacheDirXDG()
}

func CacheDirLegacy() (path string, exists bool) {
	dir := homeDir(".scaffold/cache")
	_, err := os.Stat(dir)
	if err != nil {
		return "", false
	}

	return dir, true
}

func CacheDirXDG() string {
	cacheDir := os.Getenv("XDG_DATA_HOME")
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get home directory")
		}

		cacheDir = filepath.Join(home, ".local", "share")
	}

	return filepath.Join(cacheDir, "scaffold", "templates")
}

func homeDir(s ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get home directory")
	}

	return filepath.Join(append([]string{home}, s...)...)
}

// MigrateLegacyPaths migrates the legacy configuration and cache directories to the
// new XDG Base Directory Specification paths.
func MigrateLegacyPaths() error {
	legacyFilepath, exists := RCFilepathLegacy()
	if exists {
		err := migrateRCFile(legacyFilepath)
		if err != nil {
			return err
		}
	}

	legacyFilepath, exists = CacheDirLegacy()
	if exists {
		err := migrateCacheDir(legacyFilepath)
		if err != nil {
			return err
		}
	}

	return nil
}

func migrateRCFile(legacyFilepath string) error {
	xdgFilepath := RCFilepathXDG()
	err := os.MkdirAll(filepath.Dir(xdgFilepath), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Rename(legacyFilepath, xdgFilepath)
	if err != nil {
		log.Error().Err(err).Msg("failed to migrate legacy configuration file")
		return err
	}

	log.Info().Str("from", legacyFilepath).Str("to", xdgFilepath).Msg("migrated legacy configuration file")
	return nil
}

func migrateCacheDir(legacyDir string) error {
	xdgDir := CacheDirXDG()

	err := os.Rename(legacyDir, xdgDir)
	if err != nil {
		log.Error().Err(err).Msg("failed to migrate legacy cache directory")
		return err
	}

	log.Info().Str("from", legacyDir).Str("to", xdgDir).Msg("migrated legacy cache directory")
	return nil
}
