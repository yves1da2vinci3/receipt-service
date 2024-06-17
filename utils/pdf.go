package utils

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Format represents the paper size dimensions.
type Format struct {
	Title  string
	Width  float64
	Height float64
}

// Define the available formats
var formats = []Format{
	{"A1", 23.39, 33.11},
	{"A2", 16.54, 23.39},
	{"A3", 11.69, 16.54},
	{"A4", 8.27, 11.69},
	{"A5", 5.83, 8.27},
	{"A6", 4.13, 5.83},
}

// GetFormat retrieves the format based on the title.
func GetFormat(title string) (*Format, error) {
	for _, f := range formats {
		if f.Title == title {
			return &f, nil
		}
	}
	return nil, errors.New("format not found")
}

// GeneratePDF generates a PDF from the given HTML content and saves it to a file
func GeneratePDF(htmlContent string, format string, outputPath string) ([]byte, error) {
	// Ensure the directory exists
	dir := filepath.Dir(outputPath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Retrieve the format
	f, err := GetFormat(format)
	if err != nil {
		return nil, err
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	err = chromedp.Run(ctx, printToPDF(htmlContent, &buf, f))
	if err != nil {
		return nil, err
	}

	// Save the PDF to the specified path
	err = SavePDF(outputPath, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func printToPDF(htmlContent string, res *[]byte, format *Format) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("data:text/html," + url.PathEscape(htmlContent)),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(format.Width).
				WithPaperHeight(format.Height).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

// SavePDF saves the PDF bytes to a specified file path
func SavePDF(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
