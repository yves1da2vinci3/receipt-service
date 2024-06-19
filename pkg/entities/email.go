package entities

type EmailRequest struct {
	ReceiptType string `json:"receiptType"`
	Data        struct {
		LogementName    string `json:"LogementName"`
		ReservationDate string `json:"ReservationDate"`
		StartDate       string `json:"StartDate"`
		EndDate         string `json:"EndDate"`
		UserName        string `json:"UserName"`
		ReservationID   string `json:"ReservationID"`
		Email           string `json:"Email"`
	} `json:"data"`
}
