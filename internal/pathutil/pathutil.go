package pathutil

import (
	"fmt"
	"os"
)

func GetProjectPath(path string) (string, error) {
	if len(path) == 0 || path == "" {
		return "", fmt.Errorf("missing project path argument")
	}

	return path, nil
}

func ValidateDir(projectRootPath string) bool {
	fileInfo, err := os.Stat(projectRootPath)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func ValidateFile(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if fileInfo.Mode().IsRegular() {
		return true

	}
	return false
}
