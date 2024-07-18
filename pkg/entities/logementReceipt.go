package entities

type LogementReceiptData struct {
	LogementName    string  `json:"logementName"`
	ReservationDate string  `json:"reservationDate"`
	StartDate       string  `json:"startDate"`
	EndDate         string  `json:"endDate"`
	UserName        string  `json:"userName"`
	ReservationID   string  `json:"reservationId"`
	QRCode          *string `json:"qrcode,omitempty"`
}
