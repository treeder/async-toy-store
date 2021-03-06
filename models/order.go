package models

// Order ...
// todo: this struct and the nats code in main() should be generated by AsyncAPI
type Order struct {
	ID         string  `json:"id"`
	Amount     float64 `json:"amount"`
	Comment    string  `json:"comment"`
	Status     string  `json:"status"`
	PaymentID  string  `json:"payment_id"`
	TrackingID string  `json:"tracking_id"`
}
