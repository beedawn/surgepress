package help

import (
	"testing"
)

func TestUsageValid(t *testing.T) {
	args := []string{"./main", "markdownfiles/", "config.json"}
	helpString, err := Usage(args)
	if err != nil {
		t.Fatal(err)
	}
	want := `Usage:
  ./main <file_directory> <config>

Arguments:
  file_directory    Path of the markdown files to generate site 
  config 			Path to config json file

Options:
  -h, --help  Show this help message.
`
	if helpString != want {
		t.Errorf("help output = %s, want %s", helpString, want)
	}

}

func TestUsageEmpty(t *testing.T) {
	helpString := []string{}
	_, err := Usage(helpString)

	if err == nil {
		t.Errorf("help output = %v, want nil", err)
	}
}
