package service

import (
	"document-generator/pkg/gotenberg"
	"fmt"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
)

func NewDocumentService(templatesDir string, gotClient *gotenberg.Client) *DocumentService {
	return &DocumentService{templatesDir: templatesDir, gotClient: gotClient}
}

func (s *DocumentService) RenderHTML(code string, data map[string]any) (string, error) {
	tplPath := filepath.Join(s.templatesDir, "html", code+".html")
	if _, err := os.Stat(tplPath); os.IsNotExist(err) {
		return "", fmt.Errorf("render html: template %q not found in %s", code, s.templatesDir)
	}

	tpl, err := pongo2.FromFile(tplPath)
	if err != nil {
		return "", fmt.Errorf("render html: parse template %q: %w", tplPath, err)
	}

	out, err := tpl.Execute(data)
	if err != nil {
		return "", fmt.Errorf("render html: execute template %q: %w", tplPath, err)
	}

	return out, nil
}

func (s *DocumentService) RenderPDF(code string, data map[string]any, filename string) ([]byte, error) {
	html, err := s.RenderHTML(code, data)
	if err != nil {
		return nil, err
	}
	return s.gotClient.HTMLToPDF(filename, html)
}
