package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCodeFile generates a QR code image file and returns the file path
func GenerateQRCodeFile(content string) (string, error) {
	fmt.Printf("Generating QR code: %s\n", content)
	fileName := fmt.Sprintf("/qrcode_%s.png", filepath.Base(content))
	filePath := filepath.Join("uploads", "qrcodes", fileName)

	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}
	err = qrcode.WriteFile(content, qrcode.Medium, 256, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// DeleteQRCodeFile deletes the specified QR code image file
func DeleteQRCodeFile(filePath string) error {
	return os.Remove(filePath)
}

type QrcodeManager struct {
	content string
	path    string
}

func (qm *QrcodeManager) SetContent(content string) {
	qm.content = content
}

func (qm *QrcodeManager) GetContent() string {
	return qm.content
}

func (qm *QrcodeManager) SetPath(path string) {
	qm.path = path
}

func (qm *QrcodeManager) GetPath() string {
	return qm.path
}
