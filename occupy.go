package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type occupy struct {
	ID             int    `json:"id"`
	CustomerID     int    `json:"customer_id"`
	TableID        int    `json:"table_id"`
	Type           string `json:"type"`
	OccupationDate string `json:"occupation_date"`
	Time           string `json:"time"`
	Status         string `json:"status"`
	GuestsNumber   int    `json:"guests_number"`
	Notes          string `json:"notes"`
	Duration       string `json:"duration"`
	WaiterID       int    `json:"waiter_id"`
	AddDate        string `json:"add_date"`
	CustomerName   string `json:"customer_name"`
	CustomerPhone  string `json:"customer_phone"`
	TableMaxChairs int    `json:"table_max_chairs"`
	TableNumber    int    `json:"table_number"`
}

func getReservations(db *sql.DB, restaurantID, start, count int) ([]occupy, error) {
	statement := fmt.Sprintf("SELECT occupy.id, occupy.type, occupy.time, occupy.occupation_date, occupy.status, occupy.guests_number, occupy.notes, case when customers.name is Null then '' else customers.name end , case when customers.phone is Null then '' else customers.phone end, tables.max_chairs,tables.table_number FROM occupy join tables on occupy.table_id = tables.id left join customers on customers.id = occupy.customer_id join sections on sections.id = tables.section_id join floors on floors.id = sections.floor_id join restaurants on restaurants.id = floors.restaurant_id where restaurants.id = %d LIMIT %d OFFSET %d", restaurantID, count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	occuppy := []occupy{}
	for rows.Next() {
		var o occupy
		if err := rows.Scan(&o.ID, &o.Type, &o.Time, &o.OccupationDate, &o.Status, &o.GuestsNumber, &o.Notes, &o.CustomerName, &o.CustomerPhone, &o.TableMaxChairs, &o.TableNumber); err != nil {
			return nil, err
		}
		occuppy = append(occuppy, o)
	}
	return occuppy, nil
}

func (o *occupy) addReservation(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO `occupy`(`customer_id`, `table_id`, `type`, `occupation_date`, `time`, `status`, `guests_number`, `notes`, `duration`, `waiter_id`, `add_date`) VALUES (%d,%d,'%s','%s','%s','%s',%d,'%s','%s',%d,'%s')", o.CustomerID, o.TableID, o.Type, o.Time, o.OccupationDate, o.Status, o.GuestsNumber, o.Notes, o.Duration, o.WaiterID, time.Now())
	log.Println(statement)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

func (o *occupy) deleteReservation(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM `occupy` WHERE id=%d", o.ID)
	_, err := db.Exec(statement)
	return err
}
