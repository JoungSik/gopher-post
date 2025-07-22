package email

import (
	"fmt"
	"strconv"

	"github.com/JoungSik/gopher-post/internal/config"
	"github.com/JoungSik/gopher-post/internal/feed"
	"github.com/JoungSik/gopher-post/internal/template"
	"gopkg.in/mail.v2"
)

type Service struct {
	config          *config.SMTPConfig
	templateService *template.Service
}

func NewService(config *config.SMTPConfig, templateService *template.Service) *Service {
	return &Service{
		config:          config,
		templateService: templateService,
	}
}

func (s *Service) SendNewsletter(recipient string, articles []feed.Article) error {
	if len(articles) == 0 {
		return fmt.Errorf("no articles to send")
	}

	subject := fmt.Sprintf("📬 고퍼 포스트 - 오늘 배달된 글 %d개", len(articles))

	body, err := s.templateService.RenderNewsletter(articles)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	return s.sendEmail(recipient, subject, body)
}

func (s *Service) sendEmail(to, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Gopher Post <%s>", s.config.FromEmail))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	port, err := strconv.Atoi(s.config.Port)
	if err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}

	d := mail.NewDialer(s.config.Host, port, s.config.Username, s.config.Password)

	return d.DialAndSend(m)
}
