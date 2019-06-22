package models

type customers struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Phone     string `json:"phone"`
	BirthDate string `json:"birth_date"`
	Company   string `json:"company"`
	JobTitle  string `json:"job_title"`
	Status    string `json:"status"`
	Sex       string `json:"sex"`
	AddDate   int    `json:"add_date"`
}
