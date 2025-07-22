package main

import (
	"log"
	"os"

	"github.com/JoungSik/gopher-post/internal/config"
	"github.com/JoungSik/gopher-post/internal/email"
	"github.com/JoungSik/gopher-post/internal/feed"
	"github.com/JoungSik/gopher-post/internal/template"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Gopher Post...")

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading: %v", err)
	}

	// Load configuration
	feedsConfig, err := config.LoadFeeds("feeds.yml")
	if err != nil {
		log.Fatalf("Error loading feeds config: %v", err)
	}

	recipientsConfig, err := config.LoadRecipients("recipients.yml")
	if err != nil {
		log.Fatalf("Error loading recipients config: %v", err)
	}

	smtpConfig := config.LoadSMTPConfig()
	if smtpConfig.FromEmail == "" || smtpConfig.Username == "" || smtpConfig.Password == "" {
		log.Fatal("SMTP configuration is incomplete. Please set SMTP_USERNAME, SMTP_PASSWORD, and FROM_EMAIL environment variables.")
	}

	// Initialize services
	parser := feed.NewParser()

	templateService, err := template.NewService("templates")
	if err != nil {
		log.Fatalf("Error initializing template service: %v", err)
	}

	emailService := email.NewService(smtpConfig, templateService)

	// Parse RSS feeds
	log.Printf("Parsing %d feeds...", len(feedsConfig.Feeds))
	articles, err := parser.ParseFeeds(feedsConfig.Feeds)
	if err != nil {
		log.Fatalf("Error parsing feeds: %v", err)
	}

	if len(articles) == 0 {
		log.Println("No new articles found. Exiting.")
		os.Exit(0)
	}

	// Send emails to recipients
	log.Printf("Sending newsletter to %d recipients...", len(recipientsConfig.Recipients))
	for _, recipient := range recipientsConfig.Recipients {
		err := emailService.SendNewsletter(recipient.Email, articles)
		if err != nil {
			log.Printf("Error sending email to %s: %v", recipient.Email, err)
			continue
		}
		log.Printf("Newsletter sent successfully to %s", recipient.Email)
	}

	log.Println("Gopher Post completed successfully!")
}
