package appdata

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataDir_EnvVarOverride(t *testing.T) {
	expected := t.TempDir()
	t.Setenv(DataDirEnvVar, expected)

	path, err := DataDir()
	assert.NoError(t, err)
	assert.Equal(t, expected, path)
}

func TestDataDir_DefaultsWhenEnvVarUnset(t *testing.T) {
	t.Setenv(DataDirEnvVar, "")

	path, err := DataDir()
	assert.NoError(t, err)
	assert.NotEmpty(t, path)

	info, err := os.Stat(path)
	assert.NoError(t, err, "expected DataDir to have created the directory")
	assert.True(t, info.IsDir())
}

func TestDataDir_DefaultPathMatchesOSConvention(t *testing.T) {
	t.Setenv(DataDirEnvVar, "")

	path, err := DataDir()
	assert.NoError(t, err)

	switch runtime.GOOS {
	case "windows":
		assert.Contains(t, path, filepath.Join("AppData", "Roaming", "zenzore"))
	case "darwin":
		assert.Contains(t, path, filepath.Join("Library", "Application Support", "zenzore"))
	default:
		assert.Contains(t, path, filepath.Join(".local", "share", "zenzore"))
	}
}

func TestDataDir_RespectsXDGDataHome(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		t.Skip("XDG_DATA_HOME only applies on linux/unix defaults")
	}

	t.Setenv(DataDirEnvVar, "")
	customXDG := t.TempDir()
	t.Setenv("XDG_DATA_HOME", customXDG)

	path, err := DataDir()
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(customXDG, "zenzore"), path)
}

func TestDataDir_CreatesDirectoryIfMissing(t *testing.T) {
	target := filepath.Join(t.TempDir(), "fresh-zenzore-data")
	t.Setenv(DataDirEnvVar, target)

	_, err := os.Stat(target)
	assert.Error(t, err, "expected target not to exist yet")

	path, err := DataDir()
	assert.NoError(t, err)
	assert.Equal(t, target, path)

	info, err := os.Stat(target)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}
