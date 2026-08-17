package argparser

import (
	"errors"
)

func ArgParser(args []string) (string, string, error) {
	if len(args) != 3 || args[1] == "-h" || args[1] == "--help" {
		return "", "", errors.New("Invalid arguments")
	}
	configFilePath := args[2]
	projectPath := args[1]
	return projectPath, configFilePath, nil
}
