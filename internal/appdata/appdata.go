package appdata

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const DataDirEnvVar = "ZENZORE_DATA_DIR"

// DataDir resolves where Zenzore should read/write its operating data.
// It checks ZENZORE_DATA_DIR first; if unset, it falls back to an
// OS-appropriate per-user data directory. Either way, the directory
// is created if it doesn't already exist before returning.
func DataDir() (string, error) {
	path := os.Getenv(DataDirEnvVar)
	if path == "" {
		var err error
		path, err = defaultDataDir()
		if err != nil {
			return "", err
		}
	}

	if err := os.MkdirAll(path, 0o755); err != nil {
		return "", fmt.Errorf("creating data dir %q: %w", path, err)
	}

	return path, nil
}

// defaultDataDir computes the OS-appropriate per-user data directory,
// without creating it. Used by DataDir when ZENZORE_DATA_DIR is unset.
func defaultDataDir() (string, error) {
	base, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving user home dir: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(base, "AppData", "Roaming")
		}
		return filepath.Join(appData, "zenzore"), nil
	case "darwin":
		return filepath.Join(base, "Library", "Application Support", "zenzore"), nil
	default: // linux and other unix-likes
		xdgData := os.Getenv("XDG_DATA_HOME")
		if xdgData == "" {
			xdgData = filepath.Join(base, ".local", "share")
		}
		return filepath.Join(xdgData, "zenzore"), nil
	}
}
