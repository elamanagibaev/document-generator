package service

import "document-generator/pkg/gotenberg"

type DocumentService struct {
	templatesDir string
	gotClient    *gotenberg.Client
}

type DocxService struct {
	templatesDir string
}

type ExcelService struct {
	templatesDir string
}

type GeneratorService struct {
	docService  *DocumentService
	xlsService  *ExcelService
	docxService *DocxService
	gotClient   *gotenberg.Client
}

type GenerateRequest struct {
	Code     string                 `json:"code"`
	Format   string                 `json:"format"`
	Filename string                 `json:"filename"`
	Data     map[string]interface{} `json:"data"`
}
