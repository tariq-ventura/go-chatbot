package books_app

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/tariq-ventura/go-chatbot/internal/logs"
)

func ExtractTextFromPDF(file multipart.File) (string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(data)

	pdfReader, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", err
	}

	var textBuilder strings.Builder
	totalPages := pdfReader.NumPage()

	logs.LogInfo("Processing PDF", map[string]any{"total_pages": totalPages})

	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		page := pdfReader.Page(pageIndex)
		if page.V.IsNull() {
			logs.LogInfo("Skipping null page", map[string]any{"page": pageIndex})
			continue
		}

		txt, err := page.GetPlainText(nil)
		if err != nil {
			logs.LogError("Error extracting text from page", map[string]any{
				"page":  pageIndex,
				"error": err.Error(),
			})
			continue
		}

		logs.LogInfo("Extracted text from page", map[string]any{
			"page":        pageIndex,
			"text_length": len(txt),
		})

		textBuilder.WriteString(txt)
	}

	extractedText := textBuilder.String()

	if len(extractedText) == 0 {
		return "", fmt.Errorf("no text could be extracted from PDF (might be image-based or protected)")
	}

	logs.LogInfo("PDF text extraction completed", map[string]any{
		"total_text_length": len(extractedText),
	})

	return extractedText, nil
}
