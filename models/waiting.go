package models

type waiting struct {
	ID           int    `json:"id"`
	CustomerID   int    `json:"customer_id"`
	Time         int    `json:"time"`
	GuestsNumber int    `json:"guests_number"`
	Notes        string `json:"notes"`
	AddDate      int    `json:"add_date"`
}
