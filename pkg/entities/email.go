package entities

type EmailRequest struct {
	ReceiptType string      `json:"receiptType"`
	Data        interface{} `json:"data"`
	UserEmail   string      `json:"email"`
}
