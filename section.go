package main

import (
	"database/sql"
	"fmt"
)

type section struct {
	ID            int    `json:"id"`
	FloorID       int    `json:"floor_id"`
	SectionNumber int    `json:"section_number"`
	SectionName   string `json:"section_name"`
	AddDate       string `json:"add_date"`
}

func (s *section) getSection(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT sections.id,sections.floor_id, sections.section_number, sections.section_name, sections.add_date FROM sections where sections.floor_id = %d", s.ID)
	return db.QueryRow(statement).Scan(&s.ID, &s.FloorID, &s.SectionNumber, &s.SectionName, &s.AddDate)
}

func getSections(db *sql.DB, floorID int) ([]section, error) {
	statement := fmt.Sprintf("SELECT sections.id,sections.floor_id, sections.section_number, sections.section_name, sections.add_date FROM sections where sections.floor_id = %d", floorID)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Section := []section{}
	for rows.Next() {
		var s section
		if err := rows.Scan(&s.ID, &s.FloorID, &s.SectionNumber, &s.SectionName, &s.AddDate); err != nil {
			return nil, err
		}
		Section = append(Section, s)
	}
	return Section, nil
}
