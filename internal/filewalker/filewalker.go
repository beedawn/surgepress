package filewalker

import (
	"fmt"
	"github.com/beedawn/surgepress/internal/content"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func WalkFiles(project_root_path string) ([]content.Page, error) {

	var pages []content.Page
	if project_root_path == "" || len(project_root_path) == 0 {
		return nil, fmt.Errorf("path is blank")
	}
	fileInfo, err := os.Stat(project_root_path)
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("expected a directory")
	}
	err = filepath.WalkDir(project_root_path, func(fileName string, d fs.DirEntry, err error) error {
		//process markdown files, do we want to look at a certain dir maybe?
		//not reached by test, do we need?
		if err != nil {
			return err
		}
		//skips processing directories
		if skipDir(d) {
			return nil
		}
		//not reached by test do we need?
		if filepath.Ext(fileName) != ".md" {
			return nil
		}

		relativePath, _ := filepath.Rel(project_root_path, fileName)

		//for now process data and write contents to html files in an out dir?
		trunc_name := strings.TrimSuffix(relativePath, filepath.Ext(relativePath))
		outPath := filepath.Join("out", trunc_name+".html")

		page := content.Page{
			SourcePath: fileName,
			OutputPath: outPath,
		}

		pages = append(pages, page)

		return nil
	})

	if err != nil {
		return nil, err
	}
	return pages, nil
}

func skipDir(d fs.DirEntry) bool {
	return d.IsDir()
}
