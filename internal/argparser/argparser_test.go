package argparser

import "testing"

func TestArgParserValid(t *testing.T) {
	args := []string{"./main", "markdownfiles/", "config.json"}

	projectPath, configPath, err := ArgParser(args)
	if err != nil {
		t.Fatalf("ArgParser returned unexpected error: %v", err)
	}

	if projectPath != args[1] {
		t.Errorf("projectPath = %q, want %q", projectPath, args[1])
	}

	if configPath != args[2] {
		t.Errorf("configPath = %q, want %q", configPath, args[2])
	}
}

func TestArgParserEmpty(t *testing.T) {
	args := []string{}

	_, _, err := ArgParser(args)
	if err == nil {
		t.Fatal("ArgParser with no arguments should return an error")
	}
}
