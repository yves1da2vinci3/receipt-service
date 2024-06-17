package utils

import (
	"errors"
)

// PDFOptions defines the format options for generating the PDF
type PDFOptions struct {
	Width           float64 // Width of the page in inches
	Height          float64 // Height of the page in inches
	PrintBackground bool    // Whether to print background graphics
}

// GetDimensions returns the width and height in inches based on the format string
func GetDimensions(format string) (float64, float64, error) {
	switch format {
	case "A4":
		return 8.27, 11.69, nil
	case "Letter":
		return 8.5, 11, nil
	// Add more formats as needed
	default:
		return 0, 0, errors.New("unsupported format")
	}
}
