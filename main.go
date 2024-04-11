package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

type FeedItem struct {
	Title   string
	Link    string
	PubDate string
}

func createReadme() error {
	tplBytes, err := os.ReadFile("README.md.tpl")
	if err != nil {
		return fmt.Errorf("Something went wrong reading the README.tpl file: %w", err)
	}
	tpl := string(tplBytes)

	lastArticles, err := getLatestArticles()
	if err != nil {
		return err
	}

	return os.WriteFile("README.md", []byte(strings.ReplaceAll(tpl, "%{{latest_articles}}%", lastArticles)), 0644)
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
	for i, item := range posts[:5] {
		if i > 0 {
			stringBuilder.WriteString(" \n")
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
		fmt.Println("README.md file generated correctly")
	}
}
