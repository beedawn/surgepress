package argparser

import (
	"fmt"
	"testing"
)

func TestArgParserValid(t *testing.T) {
	helpString := []string{"./main", "markdownfiles/", "config.json"}
	projectPath, configPath, err := ArgParser(helpString)
	if projectPath != helpString[1] && configPath != helpString[2] {
		fmt.Errorf("argues not parsed properly")
	}
	if err != nil {
		fmt.Errorf("problem parsing arguments")
	}
}

func TestUsageEmpty(t *testing.T) {
	helpString := []string{}
	_, _, err := ArgParser(helpString)
	if err == nil {
		fmt.Errorf("help output = %v, want nil", err)
	}
}
