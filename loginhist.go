package main

import (
	"database/sql"
	"fmt"
	"time"
)

type loginhist struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Flag       int    `json:"flag"`
	LoginDate  int    `json:"login_date"`
	DeviceType string `json:"devicetype"`
}

func (l *loginhist) addHist(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO loginhist(user_id, flag, login_date, device_type) VALUES(%d, %d, '%s', '%s')", l.UserID, l.Flag, time.Now(), l.DeviceType)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}
