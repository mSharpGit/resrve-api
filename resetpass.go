package main

import (
	"database/sql"
	"fmt"
	"log"
)

type resetpass struct {
	ID        int    `json:"id"`
	UserID    int    `json:"users_id"`
	Confirmed int    `json:"confirmed"`
	ResetDate string `json:"reset_date"`
	Code      string `json:"code"`
	OldPass   string `json:"old_pass"`
	NewPass   string `json:"new_pass"`
}

func (r *resetpass) addResetPass(db *sql.DB) error {

	statement := fmt.Sprintf("INSERT INTO resetpass (users_id, code, confirmed, old_pass, reset_date) VALUES (%d, '%s', %d, '%s', '%s')", r.UserID, r.Code, r.Confirmed, r.OldPass, r.ResetDate)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&r.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *resetpass) updateNewPass(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE resetpass set new_pass = '%s', confirmed = %d where id = %d and confirmed = %d", r.NewPass, r.Confirmed, r.ID, 0)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (r *resetpass) getResetPass(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, users_id, confirmed, reset_date, code, old_pass, new_pass from resetpass where id = %d", r.ID)
	log.Println(statement)
	err := db.QueryRow(statement).Scan(&r.ID, &r.UserID, &r.Confirmed, &r.ResetDate, &r.Code, &r.OldPass, &r.NewPass)
	if err != nil {
		return err
	}
	return nil
}
