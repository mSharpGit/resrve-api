package models

type bills struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	Amount     int    `json:"age"`
	Notes      string `json:"notes"`
	AddDate    int    `json:"add_date"`
}
