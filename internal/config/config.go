package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Feed struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	RSS  string `yaml:"rss"`
}

type FeedsConfig struct {
	Feeds []Feed `yaml:"feeds"`
}

type Recipient struct {
	Email string `yaml:"email"`
}

type RecipientsConfig struct {
	Recipients []Recipient `yaml:"recipients"`
}

type SMTPConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	FromEmail string
}

func LoadFeeds(filename string) (*FeedsConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config FeedsConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadRecipients(filename string) (*RecipientsConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config RecipientsConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		Host:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		Port:      getEnv("SMTP_PORT", "587"),
		Username:  getEnv("SMTP_USERNAME", ""),
		Password:  getEnv("SMTP_PASSWORD", ""),
		FromEmail: getEnv("FROM_EMAIL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
