package indexpost

import "html/template"

type IndexPost struct {
	Title   string
	Content template.HTML
	Date    string
	URL     string
}
