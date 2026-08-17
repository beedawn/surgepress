package pagebuilder

import (
	"bytes"
	"fmt"
	"github.com/beedawn/surgepress/internal/content"
	"github.com/beedawn/surgepress/internal/filewriter"
	"github.com/beedawn/surgepress/internal/indexpost"
	"github.com/beedawn/surgepress/internal/rssbuilder"
	"github.com/beedawn/surgepress/internal/siteconfig"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"html/template"
	"os"
	"path/filepath"
)

//goldmark is markdown parser...

func BuildPages(pages []content.Page, templateDir string, configData siteconfig.MetaData) error {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	//needs to be out here outside of loop or it gets added every loop
	globalStylesheets := addPathToStylesheets("css", configData.Global_Stylesheets)
	for _, page := range pages {
		pageData, err := os.ReadFile(page.SourcePath)
		if err != nil {
			return err
		}

		var htmlContent bytes.Buffer
		context := parser.NewContext()
		var mdContent bytes.Buffer
		if err := md.Convert(pageData, &mdContent, parser.WithContext(context)); err != nil {
			return err
		}

		metaData := meta.Get(context)

		data := map[string]interface{}{
			"Content": template.HTML(mdContent.Bytes()),
		}

		for key, value := range metaData {
			//adds path to beginning of stylesheet
			if key == "Stylesheets" {
				if stylesheets, ok := value.([]interface{}); ok {
					strs := make([]string, len(stylesheets))
					for i, v := range stylesheets {
						strs[i] = v.(string)
					}
					value = addPathToStylesheets("css", strs)
				}

			}

			data[key] = value
		}
		//should maybe move this path adding stuff elsewhere, but this works for now...

		data["GlobalStylesheets"] = globalStylesheets
		fmt.Println("pages")
		fmt.Println(data)
		templateName, _ := metaData["Template"].(string)
		if templateName == "" {
			templateName = "default.html"
		}

		joinedTemplatePath := filepath.Join(templateDir, templateName)

		pageTmpl, err := template.ParseFiles(joinedTemplatePath)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", joinedTemplatePath, err)
		}
		if err := pageTmpl.ExecuteTemplate(&htmlContent, templateName, data); err != nil {
			return err
		}
		renderedPage := &content.RenderedPage{
			OutputPath: page.OutputPath,
			HTML:       htmlContent.Bytes(),
		}
		if err := filewriter.WriteFile(renderedPage); err != nil {
			return fmt.Errorf("failed to write page %s: %w", page.OutputPath, err)
		}
	}
	return nil
}

func BuildIndex(pages []content.Page, templateDir string, configData siteconfig.MetaData) error {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	//var renderedPages []renderedPage

	var posts []indexpost.IndexPost
	for _, page := range pages {
		pageData, err := os.ReadFile(page.SourcePath)
		if err != nil {
			return err
		}

		context := parser.NewContext()
		var mdContent bytes.Buffer
		if err := md.Convert(pageData, &mdContent, parser.WithContext(context)); err != nil {
			return err
		}

		metaData := meta.Get(context)
		title, ok := metaData["Title"].(string)
		if ok == false {
			return fmt.Errorf("title not found in front matter of %v", page.SourcePath)
		}
		if title == "" {
			return fmt.Errorf("title is empty")
		}
		date, ok := metaData["Date"].(string)
		if ok == false {
			return fmt.Errorf("date is not in frontmatter of %v", page.SourcePath)
		}
		if date == "" {
			return fmt.Errorf("date is empty")
		}
		posts = append(posts, indexpost.IndexPost{
			Title:   title,
			Content: template.HTML(mdContent.Bytes()),
			Date:    date,
			URL:     page.OutputPath,
		})
	}

	data := map[string]interface{}{
		"Posts":             posts,
		"GlobalStylesheets": configData.Global_Stylesheets, // From config file

	}

	// send title, date and content to RSS builder
	if err := rssbuilder.RSSGenerator(posts, configData); err != nil {
		return fmt.Errorf("failed to generate RSS feed: %w", err)
	}

	var htmlContent bytes.Buffer

	pageTmpl, err := template.ParseFiles(filepath.Join(templateDir, "index.html"))
	if err != nil {
		return err
	}
	if err := pageTmpl.ExecuteTemplate(&htmlContent, "index.html", data); err != nil {
		return err
	}

	renderedPage := &content.RenderedPage{
		OutputPath: "out/index.html",
		HTML:       htmlContent.Bytes(),
	}
	return filewriter.WriteFile(renderedPage)

}

func addPathToStylesheets(path string, stylesheets []string) []string {
	for i, _ := range stylesheets {
		stylesheets[i] = filepath.Join(path, stylesheets[i])
	}
	return stylesheets
}
