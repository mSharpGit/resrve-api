package models

type restaurants struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Phone    string `json:"phone"`
	ParentID int    `json:"parent_id"`
	AddDate  int    `json:"add_date"`
}
