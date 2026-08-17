package rssbuilder

import (
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"html"
	"os"
	"surgepress/internal/indexpost"
	"surgepress/internal/siteconfig"
	"time"
)

type RSSPost struct {
	Title   string
	Content string
	Date    string
	URL     string
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *Channel `xml:"channel"`
}
type AtomLink struct {
	XMLName xml.Name `xml:"atom:link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	AtomLink    *AtomLink `xml:"atom:link,omitempty"`
	XMLName     xml.Name  `xml:"channel"`
	Items       []Item    `xml:"item"`
}
type Item struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	XMLName     xml.Name `xml:"item"`
	GUID        GUID     `xml:"guid"`
}
type GUID struct {
	Value       string `xml:",chardata"`
	IsPermaLink string `xml:"isPermaLink,attr"`
}

//TODO
//break up RRS builder into
//WriteRSS(io.Writer) → flexible core function (best design)
//WriteRSSFile(...) → convenience helper on top

func RSSGenerator(posts []indexpost.IndexPost, indexMeta siteconfig.MetaData) error {
	feed, err := RSSBuilder(posts, indexMeta)
	err = CheckOutDir(feed)
	if err != nil {
		return err
	}
	err = WriteRSS(feed, indexMeta, "out/feed.xml")
	if err != nil {
		return err
	}
	return nil
}

func RSSBuilder(posts []indexpost.IndexPost, indexMeta siteconfig.MetaData) (RSS, error) {

	var rssPosts []RSSPost
	for _, post := range posts {
		content := string(post.Content)
		//content = strings.ReplaceAll(content, "<p>", "")
		//content = strings.ReplaceAll(content, "</p>", "")
		//content = strings.ReplaceAll(content, "\n", "")
		//content = strings.ReplaceAll(content, "\r", "")
		content = html.EscapeString(content)
		rssPosts = append(rssPosts, RSSPost{
			Title:   post.Title,
			Content: content,
			Date:    post.Date,
			URL:     post.URL,
		})
	}

	feed := RSS{
		Version: "2.0",
		Channel: &Channel{
			Title:       indexMeta.Title,
			Link:        indexMeta.Link,
			Description: indexMeta.Description,
			Language:    indexMeta.Language,
			Items:       []Item{},
		},
	}
	for _, post := range rssPosts {
		//make guid for each item
		v7ID, err := uuid.NewV7()
		if err != nil {
			return RSS{}, fmt.Errorf("Failed to generate V7 UUID: %v", err)
		}
		//make publish date for each item
		var pubDate string
		t, err := time.Parse("2006-01-02", post.Date)
		if err != nil {
			pubDate = time.Now().UTC().Format(time.RFC1123Z)
		} else {
			pubDate = t.UTC().Format(time.RFC1123Z)
		}

		feed.Channel.Items = append(feed.Channel.Items, Item{
			Title:       post.Title,
			Link:        indexMeta.Link + post.URL,
			Description: post.Content,
			PubDate:     pubDate,
			GUID: GUID{Value: v7ID.String(),
				IsPermaLink: "false"},
		})

	}
	return feed, nil
}
func CheckOutDir(feed RSS) error {
	if err := os.MkdirAll("out", 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}
	return nil
}

func WriteRSS(feed RSS, indexMeta siteconfig.MetaData, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error creating output file: %v", err)
	}

	defer file.Close()

	if _, err := file.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("Error writing XML header: %v\n", err)
	}

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	if feed.Channel != nil {
		feed.Channel.AtomLink = &AtomLink{
			Href: indexMeta.Link + "/feed.xml",
			Rel:  "self",
			Type: "application/rss+xml",
		}
	}

	rssWithNamespace := struct {
		XMLName xml.Name `xml:"rss"`
		Version string   `xml:"version,attr"`
		AtomNS  string   `xml:"xmlns:atom,attr"`
		Channel *Channel `xml:"channel"`
	}{
		XMLName: xml.Name{Local: "rss"},
		Version: "2.0",
		AtomNS:  "http://www.w3.org/2005/Atom",
		Channel: feed.Channel,
	}

	if err := encoder.Encode(rssWithNamespace); err != nil {
		return fmt.Errorf("Error generating RSS: %v", err)
	}

	return nil
}
