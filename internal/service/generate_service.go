package service

import (
	"fmt"

	"document-generator/pkg/gotenberg"
)

func NewGeneratorService(doc *DocumentService, xls *ExcelService, docx *DocxService, got *gotenberg.Client) *GeneratorService {
	return &GeneratorService{
		docService:  doc,
		xlsService:  xls,
		docxService: docx,
		gotClient:   got,
	}
}

func (s *GeneratorService) Generate(req GenerateRequest) (filename string, contentType string, data []byte, err error) {
	switch req.Format {
	case "html":
		result, err := s.docService.RenderHTML(req.Code, req.Data)
		return req.Code + ".html", "text/html; charset=utf-8", []byte(result), err

	case "pdf":
		if req.Filename == "" {
			req.Filename = "document.pdf"
		}
		result, err := s.docService.RenderPDF(req.Code, req.Data, req.Filename)
		return req.Filename, "application/pdf", result, err

	case "xlsx":
		if req.Filename == "" {
			req.Filename = req.Code + ".xlsx"
		}
		result, err := s.xlsService.RenderExcel(req.Code, req.Data)
		return req.Filename, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result, err

	case "docx":
		if req.Filename == "" {
			req.Filename = req.Code + ".docx"
		}
		result, err := s.docxService.RenderDOCX(req.Code, req.Data)
		return req.Filename, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", result, err

	default:
		return "", "", nil, fmt.Errorf("unsupported format: %s", req.Format)
	}
}
