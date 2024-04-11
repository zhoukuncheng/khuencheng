package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedItem struct {
	Title   string
	Link    string
	PubDate string
}

// Add a new struct to hold template data
type ReadmeData struct {
	LatestArticles string
}

func createReadme() error {
	// Read and parse the template file
	tplBytes, err := os.ReadFile("README.md.tpl")
	if err != nil {
		return fmt.Errorf("something went wrong reading the README.tpl file: %w", err)
	}
	tpl, err := template.New("readme").Parse(string(tplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	lastArticles, err := getLatestArticles()
	if err != nil {
		return err
	}

	var readmeBuffer bytes.Buffer
	err = tpl.Execute(&readmeBuffer, ReadmeData{LatestArticles: lastArticles})
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write the filled template to README.md
	return os.WriteFile("README.md", readmeBuffer.Bytes(), 0644)
}

func getLatestArticles() (string, error) {
	posts, err := getBlogRSS()
	if err != nil {
		return "", err
	}

	// Sorting posts by publication date
	sort.Slice(posts, func(i, j int) bool {
		timeA, _ := time.Parse(time.RFC1123Z, posts[i].PubDate)
		timeB, _ := time.Parse(time.RFC1123Z, posts[j].PubDate)
		return timeA.After(timeB)
	})

	var stringBuilder strings.Builder
	for i, item := range posts[:6] { // Adjusted to safely handle less than 6 articles
		if i > 0 {
			stringBuilder.WriteString("\n")
		}
		stringBuilder.WriteString(fmt.Sprintf("* [%s](%s)", item.Title, item.Link))
	}

	return stringBuilder.String(), nil
}

func getBlogRSS() ([]FeedItem, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://zhoukuncheng.github.io/index.xml")
	if err != nil {
		return nil, err
	}

	var items []FeedItem
	for _, item := range feed.Items {
		items = append(items, FeedItem{
			Title:   item.Title,
			Link:    item.Link,
			PubDate: item.Published,
		})
	}

	return items, nil
}

func main() {
	if err := createReadme(); err != nil {
		fmt.Printf("Oops! There was an error: %v\n", err)
	} else {
		fmt.Println("README.md file generated correctly.")
	}
}
