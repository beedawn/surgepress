package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"surgepress/internal/argparser"
	"surgepress/internal/content"
	"surgepress/internal/filewalker"
	"surgepress/internal/help"
	"surgepress/internal/pagebuilder"
	"surgepress/internal/pathutil"
	"surgepress/internal/siteconfig"
)

func main() {
	projectPath, configFilePath, err := argparser.ArgParser(os.Args)
	if err != nil {
		helpString, err := help.Usage(os.Args)
		if err == nil {
			panic(err)
		}
		fmt.Fprintf(os.Stderr, "%s", helpString)
	}
	projectRootPath, err := pathutil.GetProjectPath(projectPath)
	if err != nil {
		panic(err)
	}

	templateDir := "internal/template"
	//check config file path is valid filetype
	if !pathutil.ValidateFile(configFilePath) {
		panic("Invalid config file path")
	}
	configFile, err := os.Open(configFilePath)
	defer configFile.Close()
	if !pathutil.ValidateDir(projectRootPath) {
		panic("Invalid project root path")
	}

	// 2. Initialize the target variable
	var configData siteconfig.MetaData
	//indexMeta := siteconfig.IndexMeta{
	//	Title:       "My Blog",
	//	Description: "My brand new exciting blog",
	//	Link:        "https://blog.com",
	//	Language:    "en-US",
	//}
	// 3. Create a decoder and decode the stream
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&configData); err != nil {
		log.Fatalf("failed to decode JSON: %v", err)
	}

	pages, err := filewalker.WalkFiles(projectRootPath)
	if err != nil {
		panic(err)
	}

	//also template injection? maybe as part of build pages
	//style sheets too
	err = pagebuilder.BuildPages(pages, templateDir, configData)
	if err != nil {
		panic(err)
	}

	err = pagebuilder.BuildIndex(pages, templateDir, configData)
	if err != nil {
		panic(err)
	}
}

func PrintHTML(builtPages []content.RenderedPage) {
	for i, page := range builtPages {
		fmt.Printf("Page %d:\n", i)
		fmt.Printf("  OutputPath: %s\n", page.OutputPath)
		fmt.Printf("  HTML length: %d bytes\n", len(page.HTML))
		fmt.Printf("  HTML preview: %s\n", string(page.HTML)) // Full HTML
		fmt.Println("---")
	}

}
