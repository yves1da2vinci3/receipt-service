package email

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

// EmailService interface
type EmailService interface {
	SendEmail(to string, pdfPath string) error
}

// SMTPEmailService struct
type SMTPEmailService struct {
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	SenderPasswd string
}

// NewSMTPEmailService function to create a new SMTPEmailService
func NewSMTPEmailService(host, port, email, password string) *SMTPEmailService {
	return &SMTPEmailService{
		SMTPHost:     host,
		SMTPPort:     port,
		SenderEmail:  email,
		SenderPasswd: password,
	}
}

// SendEmail method to send an email with an attached PDF
func (s *SMTPEmailService) SendEmail(to string, pdfPath string) error {
	// Read the PDF file
	pdfData, err := os.ReadFile(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to read PDF file: %v", err)
	}

	// Prepare email
	subject := "Your PDF Receipt"
	body := "Please find the attached PDF receipt."
	filename := filepath.Base(pdfPath)
	boundary := "my-boundary-12345"

	// Build the email headers
	headers := []string{
		"From: " + s.SenderEmail,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: multipart/mixed; boundary=" + boundary,
	}

	// Build the email body
	message := strings.Join(headers, "\r\n") + "\r\n\r\n" +
		"--" + boundary + "\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"Content-Transfer-Encoding: 7bit\r\n\r\n" +
		body + "\r\n\r\n" +
		"--" + boundary + "\r\n" +
		"Content-Type: application/pdf\r\n" +
		"Content-Disposition: attachment; filename=\"" + filename + "\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n\r\n" +
		encodeBase64(pdfData) + "\r\n" +
		"--" + boundary + "--"

	// Send the email
	auth := smtp.PlainAuth("", s.SenderEmail, s.SenderPasswd, s.SMTPHost)
	err = smtp.SendMail(s.SMTPHost+":"+s.SMTPPort, auth, s.SenderEmail, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// encodeBase64 function to encode data to base64 format
func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
