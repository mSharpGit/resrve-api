package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type table struct {
	ID      int `json:"id"`
	FloorID int `json:"floor_id"`
	/* SectionID   int    `json:"section_id"` */
	TableNumber int     `json:"table_number"`
	TableName   string  `json:"table_name"`
	Shape       string  `json:"shape"`
	Height      float64 `json:"height"`
	Width       float64 `json:"width"`
	Diameter    float64 `json:"diameter"`
	Xloc        float64 `json:"xloc"`
	Yloc        float64 `json:"yloc"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	MinChairs   int     `json:"min_chairs"`
	MaxChairs   int     `json:"max_chairs"`
	AddDate     string  `json:"add_date"`
}

func (t *table) getTable(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT `id`, `floor_id`, `table_number`, `table_name`, `shape`, `height`, `width`, `diameter`, `xloc`, `yloc`, `start_date`, `end_date`, `min_chairs`, `max_chairs`, `add_date` FROM `tables` WHERE id = %d", t.ID)
	return db.QueryRow(statement).Scan(&t.ID, &t.FloorID, &t.TableNumber, &t.TableName, &t.Shape, &t.Height, &t.Width, &t.Diameter, &t.Xloc, &t.Yloc, &t.StartDate, &t.EndDate, &t.MinChairs, &t.MaxChairs, &t.AddDate)
}

func (t *table) addTable(db *sql.DB) error {

	statement := fmt.Sprintf("INSERT INTO `tables`  (`floor_id`, `table_number`, `table_name`, `shape`, `height`, `width`, `diameter`, `xloc`, `yloc`, `start_date`, `end_date`, `min_chairs`, `max_chairs`, `add_date`) VALUES (%d,%d,'%s','%s',%f,%f,%f,%f,%f,'%s','%s',%d,%d,'%s')", t.FloorID, t.TableNumber, t.TableName, t.Shape, t.Height, t.Width, t.Diameter, t.Xloc, t.Yloc, t.StartDate, t.EndDate, t.MinChairs, t.MaxChairs, time.Now())
	log.Println(statement)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&t.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) editTable(db *sql.DB) error {

	statement := fmt.Sprintf("UPDATE `tables` SET `table_number` = %d,`table_name` = '%s', `shape` = '%s', `height` = %f, `width` = %f, `diameter` = %f, `xloc` = %f, `yloc` = %f, `start_date` = '%s', `end_date` = '%s', `min_chairs` = %d, `max_chairs` = %d WHERE id = %d", t.TableNumber, t.TableName, t.Shape, t.Height, t.Width, t.Diameter, t.Xloc, t.Yloc, t.StartDate, t.EndDate, t.MinChairs, t.MaxChairs, t.ID)
	log.Println(statement)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

func editTableBatch(db *sql.DB, r []table) ([]table, error) {

	for _, row := range r {
		statement := fmt.Sprintf("UPDATE `tables` SET `table_number` = %d,`table_name` = '%s', `shape` = '%s', `height` = %f, `width` = %f, `diameter` = %f, `xloc` = %f, `yloc` = %f, `start_date` = '%s', `end_date` = '%s', `min_chairs` = %d, `max_chairs` = %d WHERE id = %d", row.TableNumber, row.TableName, row.Shape, row.Height, row.Width, row.Diameter, row.Xloc, row.Yloc, row.StartDate, row.EndDate, row.MinChairs, row.MaxChairs, row.ID)
		log.Println(statement)
		_, err := db.Exec(statement)
		if err != nil {
			return nil, err
		}
	}
	//trim the last ,
	//=sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	//=stmt, _ := db.Prepare(sqlStr)

	//log.Println(sqlStr)

	//format all vals at once
	//=_, err := stmt.Exec(vals...)
	//=if err != nil {
	//=	return nil, err
	//=}

	return nil, nil
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

func (t *table) deleteTable(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM `tables` WHERE id=%d", t.ID)
	_, err := db.Exec(statement)
	return err
}
