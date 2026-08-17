package help

import "fmt"

func Usage(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no arguments provided, something went wrong")
	}
	return fmt.Sprintf(`Usage:
  %s <file_directory> <config>

Arguments:
  file_directory    Path of the markdown files to generate site 
  config 			Path to config json file

Options:
  -h, --help  Show this help message.
`, args[0]), nil
}
