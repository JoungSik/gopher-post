package template

import (
	"bytes"
	"html/template"
	"path/filepath"

	"github.com/JoungSik/gopher-post/internal/feed"
)

type Service struct {
	templates *template.Template
}

type NewsletterData struct {
	Articles []feed.Article
}

func NewService(templateDir string) (*Service, error) {
	templates, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, err
	}

	return &Service{
		templates: templates,
	}, nil
}

func (s *Service) RenderNewsletter(articles []feed.Article) (string, error) {
	data := NewsletterData{
		Articles: articles,
	}

	var buf bytes.Buffer
	err := s.templates.ExecuteTemplate(&buf, "newsletter.html", data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
