package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type customers struct {
	ID           int    `json:"id"`
	RestaurantID int    `json:"restaurant_id"`
	Title        string `json:"title"`
	Name         string `json:"name"`
	Lastname     string `json:"lastname"`
	Email        string `json:"email"`
	CountryCode  string `json:"country_code"`
	Phone        string `json:"phone"`
	BirthDate    string `json:"birth_date"`
	Company      string `json:"company"`
	JobTitle     string `json:"job_title"`
	Status       string `json:"status"`
	Sex          string `json:"sex"`
	Notes        string `json:"notes"`
	AddDate      string `json:"add_date"`
}

func getCustomers(db *sql.DB, restaurantID, start, count int) ([]customers, error) {
	statement := fmt.Sprintf("SELECT `id`, `restaurant_id`,`title` ,`name`,`lastname`, `email`, `country_code`, `phone`,`birth_date`,`company`,`job_title`,`status`,`sex`, `notes`, `add_date` FROM customers where restaurant_id = %d LIMIT %d OFFSET %d", restaurantID, count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	customersV := []customers{}
	for rows.Next() {
		var c customers
		if err := rows.Scan(&c.ID, &c.RestaurantID, &c.Title, &c.Name, &c.Lastname, &c.Email, &c.CountryCode, &c.Phone, &c.BirthDate, &c.Company, &c.JobTitle, &c.Status, &c.Sex, &c.Notes, &c.AddDate); err != nil {
			return nil, err
		}
		customersV = append(customersV, c)
	}
	return customersV, nil
}

func getCustomersBatch(db *sql.DB, IDS string) ([]customers, error) {
	statement := fmt.Sprintf("SELECT `id`, `restaurant_id`,`title` ,`name`,`lastname`, `email`, `country_code`,`phone`,`birth_date`,`company`,`job_title`,`status`,`sex`, `notes`, `add_date` FROM customers where id in %s", IDS)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	customersV := []customers{}
	for rows.Next() {
		var c customers
		if err := rows.Scan(&c.ID, &c.RestaurantID, &c.Title, &c.Name, &c.Lastname, &c.Email, &c.CountryCode, &c.Phone, &c.BirthDate, &c.Company, &c.JobTitle, &c.Status, &c.Sex, &c.Notes, &c.AddDate); err != nil {
			return nil, err
		}
		customersV = append(customersV, c)
	}
	return customersV, nil
}

func searchCustomers(db *sql.DB, term string, start, count int) ([]customers, error) {
	likeTerm := "%" + term + "%"
	statement := fmt.Sprintf("SELECT `id`, `restaurant_id`,`title` ,`name`,`lastname`,`email`, `country_code`, `phone`,`birth_date`,`company`,`job_title`,`status`,`sex`, `notes`, `add_date` FROM customers where name like '%s' LIMIT %d OFFSET %d", likeTerm, count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	customersV := []customers{}
	for rows.Next() {
		var c customers
		if err := rows.Scan(&c.ID, &c.RestaurantID, &c.Title, &c.Name, &c.Lastname, &c.Email, &c.CountryCode, &c.Phone, &c.BirthDate, &c.Company, &c.JobTitle, &c.Status, &c.Sex, &c.Notes, &c.AddDate); err != nil {
			return nil, err
		}
		customersV = append(customersV, c)
	}
	return customersV, nil
}

func (c *customers) addCustomer(db *sql.DB) error {

	statement := fmt.Sprintf("INSERT INTO `customers`  (`restaurant_id`,`title`, `name`, `lastname`, `email`, `country_code`, `phone`, `birth_date`, `company`, `job_title`, `status`, `sex`, `notes`, `add_date`) VALUES (%d,'%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')", c.RestaurantID, c.Title, c.Name, c.Lastname, c.Email, c.CountryCode, c.Phone, c.BirthDate, c.Company, c.JobTitle, c.Status, c.Sex, c.Notes, time.Now())
	log.Println(statement)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&c.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *customers) editCustomer(db *sql.DB) error {

	statement := fmt.Sprintf("UPDATE `customers` SET `restaurant_id` = %d,`title` = '%s', `name` = '%s', `lastname` = '%s', `email` = '%s', `phone` = '%s', `phone` = '%s', `birth_date` = '%s', `company` = '%s', `job_title` = '%s', `status` = '%s', `sex` = '%s', `notes` = '%s', `add_date` = '%s' WHERE id = %d", c.RestaurantID, c.Title, c.Name, c.Lastname, c.Email, c.CountryCode, c.Phone, c.BirthDate, c.Company, c.JobTitle, c.Status, c.Sex, c.Notes, time.Now(), c.ID)
	log.Println(statement)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

func (c *customers) deleteCustomer(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM `customers` WHERE id=%d", c.ID)
	_, err := db.Exec(statement)
	return err
}
