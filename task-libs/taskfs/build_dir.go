package taskfs

import (
	"fmt"
	"os"
	"path/filepath"
)

type BuildDir string

func (buildDir BuildDir) SubDir(subDir string) (string, error) {
	dir := filepath.Join(string(buildDir), subDir)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("missing sub directory '%s' in build directory '%s'", subDir, buildDir)
	}

	return dir, nil
}
