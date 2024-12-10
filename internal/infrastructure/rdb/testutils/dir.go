package testutils

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetProjectDir get project dir
func GetProjectDir() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return findProjectDir(path, 10)
}

// findV2Dir find v2 dir
func findProjectDir(path string, tryCount int) (string, error) {
	if _, err := os.Stat(path + "/go.mod"); os.IsNotExist(err) {
		if tryCount == 0 {
			return "", fmt.Errorf("project dir nothing")
		}
		return findProjectDir(filepath.Dir(path), tryCount-1)
	}

	return path, nil
}
