package gotenberg

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func (c *Client) HTMLToPDF(filename string, html string) ([]byte, error) {
	if filename == "" {
		filename = DefaultPDFName
	}
	return c.htmlToFile(filename, html)
}

func (c *Client) HTMLToDOCX(filename string, html string) ([]byte, error) {
	if filename == "" {
		filename = DefaultDOCXName
	}
	return c.htmlToFile(filename, html)
}

func (c *Client) htmlToFile(filename string, html string) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(FieldHTMLFile, FieldHTMLFile)
	if err != nil {
		return nil, err
	}
	_, _ = part.Write([]byte(html))

	_ = writer.WriteField(FieldResultFile, filename)

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+EndpointChromiumHTML, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gotenberg error: %s", string(b))
	}

	return io.ReadAll(resp.Body)
}
