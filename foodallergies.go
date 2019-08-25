package main

import (
	"database/sql"
	"fmt"
	"log"
)

type foodAllergies struct {
	ID       int    `json:"id"`
	FoodType string `json:"food_type"`
}

type allergiesLink struct {
	ID          int `json:"id"`
	CustomerID  int `json:"customer_id"`
	AllergiesID int `json:"allergy_id"`
}

func getFoodAllergies(db *sql.DB, start, count int) ([]foodAllergies, error) {
	statement := fmt.Sprintf("SELECT `id`, `food_type` FROM food_allergy LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	FoodAllergiesV := []foodAllergies{}
	for rows.Next() {
		var f foodAllergies
		if err := rows.Scan(&f.ID, &f.FoodType); err != nil {
			return nil, err
		}
		FoodAllergiesV = append(FoodAllergiesV, f)
	}
	return FoodAllergiesV, nil
}

func getFoodAllergyLink(db *sql.DB, id int) ([]foodAllergies, error) {
	statement := fmt.Sprintf("SELECT  food_allergy.id, food_allergy.food_type FROM food_allergy join allergy_link on allergy_link.allergy_id = food_allergy.id WHERE allergy_link.customer_id = %d", id)
	log.Println(statement)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	FoodAllergiesV := []foodAllergies{}
	for rows.Next() {
		var f foodAllergies
		if err := rows.Scan(&f.ID, &f.FoodType); err != nil {
			return nil, err
		}
		FoodAllergiesV = append(FoodAllergiesV, f)
	}
	return FoodAllergiesV, nil
}

func (f *foodAllergies) addFoodAllergies(db *sql.DB) error {

	statement := fmt.Sprintf("INSERT INTO `food_allergy`(`food_type`) VALUES ('%s')", f.FoodType)
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

/* func (f *allergiesLink) addFoodAllergyLink(db *sql.DB) error {

	statement := fmt.Sprintf("INSERT INTO `allergy_link`(`customer_id`, `allergy_id`) VALUES (%d, %d)", f.CustomerID, f.AllergiesID)
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
} */

func (f *allergiesLink) deleteCustAllergy(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM `allergy_link` WHERE customer_id=%d", f.CustomerID)
	_, err := db.Exec(statement)
	return err
}

func addFoodAllergyLink(db *sql.DB, r []allergiesLink) ([]foodAllergies, error) {

	sqlStr := "INSERT INTO `allergy_link`(`customer_id`, `allergy_id`) VALUES"
	vals := []interface{}{}

	for _, row := range r {
		sqlStr += "(?, ?),"
		vals = append(vals, row.CustomerID, row.AllergiesID)
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

	return nil, nil
}
