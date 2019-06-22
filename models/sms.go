package models

type sms struct {
	ID           int    `json:"id"`
	CustomerID   int    `json:"customer_id"`
	Title        string `json:"title"`
	MobileNumber string `json:"mobile_number"`
	Message      string `json:"message"`
	AddDate      int    `json:"add_date"`
}
