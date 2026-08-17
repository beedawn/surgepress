package main

import (
	"encoding/json"
	"fmt"
	"github.com/beedawn/surgepress/internal/argparser"
	"github.com/beedawn/surgepress/internal/filewalker"
	"github.com/beedawn/surgepress/internal/help"
	"github.com/beedawn/surgepress/internal/pagebuilder"
	"github.com/beedawn/surgepress/internal/pathutil"
	"github.com/beedawn/surgepress/internal/siteconfig"
	"os"
)

func main() {
	projectPath, configFilePath, err := argparser.ArgParser(os.Args)
	if err != nil {
		helpString, helpErr := help.Usage(os.Args)
		if helpErr != nil {
			fmt.Fprintf(os.Stderr, "failed to generate usage: %v\n", helpErr)
			os.Exit(1)
		}

		fmt.Fprint(os.Stderr, helpString)
		return
	}

	projectRootPath, err := pathutil.GetProjectPath(projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid project path: %v\n", err)
		os.Exit(1)
	}

	if !pathutil.ValidateDir(projectRootPath) {
		fmt.Fprintf(os.Stderr, "invalid project directory: %s\n", projectRootPath)
		os.Exit(1)
	}

	if !pathutil.ValidateFile(configFilePath) {
		fmt.Fprintf(os.Stderr, "invalid config file: %s\n", configFilePath)
		os.Exit(1)
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open config file: %v\n", err)
		os.Exit(1)
	}
	defer configFile.Close()

	var configData siteconfig.MetaData

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&configData); err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode config JSON: %v\n", err)
		os.Exit(1)
	}

	pages, err := filewalker.WalkFiles(projectRootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to discover Markdown files: %v\n", err)
		os.Exit(1)
	}

	templateDir := "internal/template"

	if err := pagebuilder.BuildPages(pages, templateDir, configData); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build pages: %v\n", err)
		os.Exit(1)
	}

	if err := pagebuilder.BuildIndex(pages, templateDir, configData); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build index: %v\n", err)
		os.Exit(1)
	}
}
