package models

type tables struct {
	ID          int    `json:"id"`
	SectionID   int    `json:"section_id"`
	TableNumber int    `json:"table_number"`
	TableName   string `json:"table_name"`
	Shape       string `json:"shape"`
	Length      int    `json:"length"`
	Width       int    `json:"width"`
	Diameter    int    `json:"diameter"`
	StartDate   int    `json:"start_date"`
	EndDate     int    `json:"end_date"`
	MinChairs   int    `json:"min_chairs"`
	MaxChairs   int    `json:"max_chairs"`
	AddDate     int    `json:"add_date"`
}
