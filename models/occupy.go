package models

type occupy struct {
	ID           int    `json:"id"`
	CustomerID   int    `json:"customer_id"`
	TableID      int    `json:"table_id"`
	Type         string `json:"type"`
	Time         int    `json:"time"`
	Status       string `json:"status"`
	GuestsNumber int    `json:"guests_number"`
	Notes        string `json:"notes"`
	Duration     string `json:"duration"`
	WaiterID     int    `json:"waiter_id"`
	AddDate      int    `json:"add_date"`
}
