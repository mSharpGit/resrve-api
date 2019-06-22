package models

type sections struct {
	ID            int    `json:"id"`
	FloorID       int    `json:"floor_id"`
	SectionNumber int    `json:"section_number"`
	SectionName   string `json:"section_name"`
	AddDate       int    `json:"add_date"`
}
