package feed

import (
	"fmt"
	"log"
	"time"

	"github.com/JoungSik/gopher-post/internal/config"
	"github.com/mmcdole/gofeed"
)

type Article struct {
	Title       string
	Link        string
	Description string
	Author      string
	Published   time.Time
	FeedName    string
}

type Parser struct {
	parser *gofeed.Parser
}

func NewParser() *Parser {
	return &Parser{
		parser: gofeed.NewParser(),
	}
}

func (p *Parser) ParseFeeds(feeds []config.Feed) ([]Article, error) {
	var articles []Article

	for _, feed := range feeds {
		log.Printf("Parsing feed: %s (%s)", feed.Name, feed.RSS)
		feedData, err := p.parser.ParseURL(feed.RSS)
		if err != nil {
			log.Printf("Error parsing feed %s: %v", feed.Name, err)
			continue
		}

		log.Printf("Found %d items in feed %s", len(feedData.Items), feed.Name)

		for _, item := range feedData.Items {
			var publishedTime time.Time
			if item.PublishedParsed != nil {
				publishedTime = *item.PublishedParsed
			} else if item.UpdatedParsed != nil {
				publishedTime = *item.UpdatedParsed
			} else {
				publishedTime = time.Now()
			}

			// Only include articles from the last 30 days for testing
			if time.Since(publishedTime) > 24*time.Hour {
				log.Printf("Skipping old article: %s (published: %v, age: %v)", item.Title, publishedTime, time.Since(publishedTime))
				continue
			}

			log.Printf("Including article: %s (published: %v)", item.Title, publishedTime)

			author := ""
			if item.Author != nil {
				author = item.Author.Name
			}

			description := item.Description
			if len(description) > 300 {
				description = description[:300] + "..."
			}

			article := Article{
				Title:       item.Title,
				Link:        item.Link,
				Description: description,
				Author:      author,
				Published:   publishedTime,
				FeedName:    feed.Name,
			}

			articles = append(articles, article)
		}
	}

	log.Printf("Found %d articles from %d feeds", len(articles), len(feeds))
	return articles, nil
}

func (p *Parser) GetRecentArticles(feeds []config.Feed, hours int) ([]Article, error) {
	var articles []Article

	for _, feed := range feeds {
		feedData, err := p.parser.ParseURL(feed.RSS)
		if err != nil {
			return nil, fmt.Errorf("error parsing feed %s: %w", feed.Name, err)
		}

		for _, item := range feedData.Items {
			var publishedTime time.Time
			if item.PublishedParsed != nil {
				publishedTime = *item.PublishedParsed
			} else if item.UpdatedParsed != nil {
				publishedTime = *item.UpdatedParsed
			} else {
				continue
			}

			// Only include articles from the specified hours
			if time.Since(publishedTime) > time.Duration(hours)*time.Hour {
				continue
			}

			author := ""
			if item.Author != nil {
				author = item.Author.Name
			}

			article := Article{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				Author:      author,
				Published:   publishedTime,
				FeedName:    feed.Name,
			}

			articles = append(articles, article)
		}
	}

	return articles, nil
}
