package feed

import (
	"fmt"
	"log"
	"regexp"
	"strings"
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

func extractTextFromHTML(htmlContent string) string {
	if htmlContent == "" {
		return ""
	}

	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(htmlContent, "")

	// Clean up extra whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	return text
}

func (p *Parser) ParseFeeds(feeds []config.Feed) ([]Article, error) {
	var articles []Article

	for _, feed := range feeds {
		log.Printf("Parsing RSS/Atom feed: %s (%s)", feed.Name, feed.RSS)
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

			// Only include articles from yesterday (Korean time zone)
			kst, _ := time.LoadLocation("Asia/Seoul")
			now := time.Now().In(kst)
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, kst)
			yesterday := today.AddDate(0, 0, -1)

			// Convert published time to Korean time zone for comparison
			publishedTimeKST := publishedTime.In(kst)

			if publishedTimeKST.Before(yesterday) || publishedTimeKST.After(today) {
				log.Printf("Skipping article not from yesterday KST: %s (published: %v KST)", item.Title, publishedTimeKST)
				continue
			}

			log.Printf("Including article: %s (published: %v)", item.Title, publishedTime)

			author := ""
			if item.Author != nil {
				author = item.Author.Name
			}

			description := item.Description
			// For Atom feeds, Description is often empty, so fallback to Content
			if description == "" && item.Content != "" {
				description = extractTextFromHTML(item.Content)
			}
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

			// Only include articles from yesterday (Korean time zone)
			kst, _ := time.LoadLocation("Asia/Seoul")
			now := time.Now().In(kst)
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, kst)
			yesterday := today.AddDate(0, 0, -1)

			// Convert published time to Korean time zone for comparison
			publishedTimeKST := publishedTime.In(kst)

			if publishedTimeKST.Before(yesterday) || publishedTimeKST.After(today) {
				continue
			}

			author := ""
			if item.Author != nil {
				author = item.Author.Name
			}

			// Handle description for GetRecentArticles too
			itemDescription := item.Description
			if itemDescription == "" && item.Content != "" {
				itemDescription = extractTextFromHTML(item.Content)
				if len(itemDescription) > 300 {
					itemDescription = itemDescription[:300] + "..."
				}
			}

			article := Article{
				Title:       item.Title,
				Link:        item.Link,
				Description: itemDescription,
				Author:      author,
				Published:   publishedTime,
				FeedName:    feed.Name,
			}

			articles = append(articles, article)
		}
	}

	return articles, nil
}
