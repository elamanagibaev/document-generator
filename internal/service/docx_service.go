package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nguyenthenguyen/docx"
)

func NewDocxService(templatesDir string) *DocxService {
	return &DocxService{templatesDir: templatesDir}
}

func (s *DocxService) RenderDOCX(code string, data map[string]interface{}) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working dir: %w", err)
	}

	tplPath := filepath.Join(wd, s.templatesDir, "docx", code+".docx")
	fmt.Println("Opening template:", tplPath)

	// Открываем шаблон
	r, err := docx.ReadDocxFile(tplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open template: %w", err)
	}
	defer r.Close()

	docx1 := r.Editable()

	// Подставляем данные
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		docx1.Replace(placeholder, fmt.Sprintf("%v", value), -1)
	}

	// Временный файл
	tmpFile := filepath.Join(os.TempDir(), "result.docx")

	// Сохраняем туда документ
	if err := docx1.WriteToFile(tmpFile); err != nil {
		return nil, fmt.Errorf("failed to write docx: %w", err)
	}

	// Читаем в память
	result, err := os.ReadFile(tmpFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read result file: %w", err)
	}

	// Удаляем временный файл
	_ = os.Remove(tmpFile)

	return result, nil
}
