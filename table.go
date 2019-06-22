package main

import (
	"database/sql"
	"fmt"
)

type table struct {
	ID          int    `json:"id"`
	SectionID   int    `json:"section_id"`
	TableNumber int    `json:"table_number"`
	TableName   string `json:"table_name"`
	Shape       string `json:"shape"`
	Length      int    `json:"length"`
	Width       int    `json:"width"`
	Diameter    int    `json:"diameter"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	MinChairs   int    `json:"min_chairs"`
	MaxChairs   int    `json:"max_chairs"`
	AddDate     string `json:"add_date"`
}

func getTables(db *sql.DB, sectionID int) ([]table, error) {
	statement := fmt.Sprintf("SELECT `id`, `section_id`, `table_number`, `table_name`, `shape`, `length`, `width`, `diameter`, `start_date`, `end_date`, `min_chairs`, `max_chairs`, `add_date` FROM `tables` WHERE section_id = %d", sectionID)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Table := []table{}
	for rows.Next() {
		var t table
		if err := rows.Scan(&t.ID, &t.SectionID, &t.TableNumber, &t.TableName, &t.Shape, &t.Length, &t.Width, &t.Diameter, &t.StartDate, &t.EndDate, &t.MinChairs, &t.MaxChairs, &t.AddDate); err != nil {
			return nil, err
		}
		Table = append(Table, t)
	}
	return Table, nil
}
