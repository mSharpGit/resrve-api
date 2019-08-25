package main

import (
	"database/sql"
	"fmt"
)

type floor struct {
	ID            int    `json:"id"`
	RestaurantID  int    `json:"restaurant_id"`
	FloorNumber   int    `json:"floor_number"`
	FloorName     string `json:"floor_name"`
	Length        int    `json:"length"`
	Width         int    `json:"width"`
	BackgroundPic string `json:"background_pic"`
	AddDate       string `json:"add_date"`
}

func (f *floor) getFloor(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT floors.id,floors.restaurant_id, floors.floor_number, floors.floor_name, floors.length, floors.width, floors.background_pic, floors.add_date FROM floors where floors.id = %d", f.ID)
	return db.QueryRow(statement).Scan(&f.ID, &f.RestaurantID, &f.FloorNumber, &f.FloorName, &f.Length, &f.Width, &f.BackgroundPic, &f.AddDate)
}

func getFloors(db *sql.DB, restaurantID int) ([]floor, error) {
	statement := fmt.Sprintf("SELECT floors.id,floors.restaurant_id, floors.floor_number, floors.floor_name, floors.length, floors.width, floors.background_pic, floors.add_date FROM floors where floors.restaurant_id = %d", restaurantID)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Floor := []floor{}
	for rows.Next() {
		var f floor
		if err := rows.Scan(&f.ID, &f.RestaurantID, &f.FloorNumber, &f.FloorName, &f.Length, &f.Width, &f.BackgroundPic, &f.AddDate); err != nil {
			return nil, err
		}
		Floor = append(Floor, f)
	}
	return Floor, nil
}
