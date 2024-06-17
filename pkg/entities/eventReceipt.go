package entities

type EventReceiptData struct {
	EventName     string `json:"eventName"`
	EventDate     string `json:"eventDate"`
	EventLocation string `json:"eventLocation"`
	UserName      string `json:"userName"`
	ReservationID string `json:"reservationId"`
}
