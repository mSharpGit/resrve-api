package main

import (
	"database/sql"
	"fmt"
	"log"
)

type table struct {
	ID      int `json:"id"`
	FloorID int `json:"floor_id"`
	/* SectionID   int    `json:"section_id"` */
	TableNumber int    `json:"table_number"`
	TableName   string `json:"table_name"`
	Shape       string `json:"shape"`
	Height      int    `json:"height"`
	Width       int    `json:"width"`
	Diameter    int    `json:"diameter"`
	Xloc        int    `json:"xloc"`
	Yloc        int    `json:"yloc"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	MinChairs   int    `json:"min_chairs"`
	MaxChairs   int    `json:"max_chairs"`
	AddDate     string `json:"add_date"`
}

func (t *table) getTable(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT `id`, `floor_id`, `table_number`, `table_name`, `shape`, `height`, `width`, `diameter`, `xloc`, `yloc`, `start_date`, `end_date`, `min_chairs`, `max_chairs`, `add_date` FROM `tables` WHERE id = %d", t.ID)
	return db.QueryRow(statement).Scan(&t.ID, &t.FloorID, &t.TableNumber, &t.TableName, &t.Shape, &t.Height, &t.Width, &t.Diameter, &t.Xloc, &t.Yloc, &t.StartDate, &t.EndDate, &t.MinChairs, &t.MaxChairs, &t.AddDate)
}

func getTables(db *sql.DB, sectionID int) ([]table, error) {
	statement := fmt.Sprintf("SELECT `id`, `floor_id`, `table_number`, `table_name`, `shape`, `height`, `width`, `diameter`, `xloc`, `yloc`, `start_date`, `end_date`, `min_chairs`, `max_chairs`, `add_date` FROM `tables` WHERE floor_id = %d", sectionID)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Table := []table{}
	for rows.Next() {
		var t table
		if err := rows.Scan(&t.ID, &t.FloorID, &t.TableNumber, &t.TableName, &t.Shape, &t.Height, &t.Width, &t.Diameter, &t.Xloc, &t.Yloc, &t.StartDate, &t.EndDate, &t.MinChairs, &t.MaxChairs, &t.AddDate); err != nil {
			return nil, err
		}
		Table = append(Table, t)
	}
	return Table, nil
}

func getTablesBatch(db *sql.DB, ids string) ([]table, error) {
	statement := fmt.Sprintf("SELECT `id`, `floor_id`, `table_number`, `table_name`, `shape`, `height`, `width`, `diameter`, `xloc`, `yloc`, `start_date`, `end_date`, `min_chairs`, `max_chairs`, `add_date` FROM `tables` WHERE id in %s", ids)
	log.Println("selecting tables in batch with ids: ", statement)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Table := []table{}
	for rows.Next() {
		var t table
		if err := rows.Scan(&t.ID, &t.FloorID, &t.TableNumber, &t.TableName, &t.Shape, &t.Height, &t.Width, &t.Diameter, &t.Xloc, &t.Yloc, &t.StartDate, &t.EndDate, &t.MinChairs, &t.MaxChairs, &t.AddDate); err != nil {
			return nil, err
		}
		Table = append(Table, t)
	}
	return Table, nil
}
