package main

import (
	"database/sql"
	"fmt"
	"log"
)

type foodRestrictions struct {
	ID       int    `json:"id"`
	FoodType string `json:"food_type"`
}

type restrictionsLink struct {
	ID            int `json:"id"`
	CustomerID    int `json:"customer_id"`
	RestrictionID int `json:"restriction_id"`
}

func getFoodRestrictions(db *sql.DB, start, count int) ([]foodRestrictions, error) {
	statement := fmt.Sprintf("SELECT `id`, `food_type` FROM food_restrictions LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	FoodRestictionsV := []foodRestrictions{}
	for rows.Next() {
		var f foodRestrictions
		if err := rows.Scan(&f.ID, &f.FoodType); err != nil {
			return nil, err
		}
		FoodRestictionsV = append(FoodRestictionsV, f)
	}
	return FoodRestictionsV, nil
}

func getFoodRestrictionLink(db *sql.DB, id int) ([]foodRestrictions, error) {
	statement := fmt.Sprintf("SELECT food_restrictions.id, food_restrictions.food_type FROM food_restrictions join restrictions_link on restrictions_link.restriction_id = food_restrictions.id WHERE restrictions_link.customer_id = %d", id)
	log.Println(statement)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	FoodRestrictionsV := []foodRestrictions{}
	for rows.Next() {
		var f foodRestrictions
		if err := rows.Scan(&f.ID, &f.FoodType); err != nil {
			return nil, err
		}
		FoodRestrictionsV = append(FoodRestrictionsV, f)
	}
	return FoodRestrictionsV, nil
}

func (f *foodRestrictions) addFoodRestrictions(db *sql.DB) error {

	statement := fmt.Sprintf("INSERT INTO `food_restrictions`(`food_type`) VALUES ('%s')", f.FoodType)
	log.Println(statement)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&f.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *restrictionsLink) deleteCustRestriction(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM `restrictions_link` WHERE customer_id=%d", f.CustomerID)
	_, err := db.Exec(statement)
	return err
}

func addFoodRestrictionsLink(db *sql.DB, r []restrictionsLink) ([]foodRestrictions, error) {

	sqlStr := "INSERT INTO `restrictions_link`(`customer_id`, `restriction_id`) VALUES"
	vals := []interface{}{}

	for _, row := range r {
		sqlStr += "(?, ?),"
		vals = append(vals, row.CustomerID, row.RestrictionID)
	}
	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, _ := db.Prepare(sqlStr)

	log.Println(sqlStr)

	//format all vals at once
	_, err := stmt.Exec(vals...)
	if err != nil {
		return nil, err
	}
	//statement := fmt.Sprintf("INSERT INTO `restrictions_link`(`customer_id`, `restriction_id`) VALUES (%d, %d)", r[0].CustomerID, r[0].RestrictionID)
	//log.Println(statement)
	/* _, err := db.Exec(stmt)
	if err != nil {
		return nil, err
	} */

	/* err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&f.ID)
	if err != nil {
		return nil, err
	} */
	return nil, nil
}
