package filewriter

import (
	"os"
	"path/filepath"
	"surgepress/internal/content"
)

func createDir(outputPath string) error {
	return os.MkdirAll(filepath.Dir(outputPath), 0755)

}

func WriteFile(page *content.RenderedPage) error {

	err := createDir(page.OutputPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(page.OutputPath, page.HTML, 0644)
	if err != nil {
		return err
	}

	return nil
}
