package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type users struct {
	ID           int    `json:"id"`
	RestaurantID int    `json:"Restaurant_id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Birthdate    string `json:"birthdate"`
	Sex          string `json:"sex"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	City         string `json:"city"`
	Country      string `json:"country"`
	PostalCode   string `json:"postalCode"`
	Confirmed    int    `json:"confirmed"`
	VerifyCode   string `json:"verifyCode"`
	Role         string `json:"role"`
	ManagerID    int    `json:"manager_id"`
	AlwaysLogged int    `json:"always_logged"`
	RegDate      string `json:"reg_date"`
}

func getWaiters(db *sql.DB, restaurantID int) ([]users, error) {
	statement := fmt.Sprintf("SELECT `id`, `restaurant_id`, `name`, `surname`, `birthdate`, `sex`, `password`, `email`, `address`, `city`, `country`, `postalcode`, `confirmed`, `verifycode`, `role`, `manager_id`, `always_logged`, `reg_date` FROM `users` WHERE restaurant_id = %d and role = 'waiter'", restaurantID)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Waiter := []users{}
	for rows.Next() {
		var u users
		if err := rows.Scan(&u.ID, &u.RestaurantID, &u.Name, &u.Surname, &u.Birthdate, &u.Sex, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.PostalCode, &u.Confirmed, &u.VerifyCode, &u.Role, &u.ManagerID, &u.AlwaysLogged, &u.RegDate); err != nil {
			return nil, err
		}
		Waiter = append(Waiter, u)
	}
	return Waiter, nil
}
func (u *users) authUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id,restaurant_id, name, birthdate, surname, sex, password, email, address, city, country, postalcode, confirmed, verifycode, reg_date, role, manager_id, always_logged FROM users WHERE email='%s'", u.Email)
	pass := u.Password

	err := db.QueryRow(statement).Scan(&u.ID, &u.RestaurantID, &u.Name, &u.Birthdate, &u.Surname, &u.Sex, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.PostalCode, &u.Confirmed, &u.VerifyCode, &u.RegDate, &u.Role, &u.ManagerID, &u.AlwaysLogged)
	if err != nil {
		return err
	}
	if u.Confirmed != 1 {
		err = errors.New("User Not Confimed Yet Please Check Your Email Address")
		return err
	}
	// Comparing the password with the hash
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)); err != nil {
		return err
	}
	/* //fmt.Println(err)
	statement = fmt.Sprintf("INSERT INTO loginhist(user_id, flag, login_date) VALUES(%d, %d, '%s')", u.ID, 1, time.Now())
	_, err = db.Exec(statement)
	if err != nil {
		return err
	} */

	return nil
}
func (u *users) verifyUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE users SET confirmed=1 WHERE verifycode='%s' and id=%d", u.VerifyCode, u.ID)
	_, err := db.Exec(statement)
	statement = fmt.Sprintf("SELECT id, email FROM users WHERE Verifycode='%s' and id =%d", u.VerifyCode, u.ID)
	err = db.QueryRow(statement).Scan(&u.ID, &u.Email)
	return err
}
func (u *users) resetUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, password FROM users WHERE email='%s' and confirmed = 1", u.Email)
	err := db.QueryRow(statement).Scan(&u.ID, &u.Password)
	if err != nil {
		return err
	}
	return nil
}
func (u *users) updateUserPass(db *sql.DB) error {

	statement := fmt.Sprintf("UPDATE users SET password='%s' WHERE id=%d", u.Password, u.ID)
	log.Println("Statment:", statement)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	/* statement := fmt.Sprintf("SELECT id, name, birthdate, surname, sex, password, email, address, city, country, postalcode, confirmed, verifycode, reg_date, role, manager_id, always_logged FROM users WHERE id = (select users_id from resetpass where id=%d and confirmed = 0) and confirmed = 1", id)
	err := db.QueryRow(statement).Scan(&u.ID, &u.Name, &u.Birthdate, &u.Surname, &u.Sex, &u.Password, &u.Email, &u.Address, &u.City, &u.Country, &u.PostalCode, &u.Confirmed, &u.VerifyCode, &u.RegDate, &u.Role, &u.ManagerID, &u.AlwaysLogged)
	if err != nil {
		return err
	} */
	/* u.Password = string(hash)
	statement = fmt.Sprintf("UPDATE resetpass SET confirmed=1, new_pass='%s' WHERE code='%s' and id=%d", hash, code, id)
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}
	statement = fmt.Sprintf("UPDATE users SET password='%s' WHERE id=%d", hash, u.ID)
	log.Println("Statment:", statement)
	_, err = db.Exec(statement)
	if err != nil {
		return err
	} */
	return nil
}
func (u *users) getUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name, birthdate FROM users WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.Name, &u.Birthdate)
}
func (u *users) updateUser(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE users SET name='%s', age=%d WHERE id=%d", u.Name, u.Birthdate, u.ID)
	_, err := db.Exec(statement)
	return err
}
func (u *users) deleteUser(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", u.ID)
	_, err := db.Exec(statement)
	return err
}
func (u *users) createUser(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name FROM users WHERE email='%s'", u.Email)
	err := db.QueryRow(statement).Scan(&u.Name)
	switch err {
	case sql.ErrNoRows:
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		fatal(err)
		date := "%m/%d/%Y"
		statement := fmt.Sprintf("INSERT INTO users(Restaurant_id, name, surname, birthdate, sex, password, email, address, city, country, postalcode, confirmed, verifycode, role, manager_id, always_logged,reg_date) VALUES(%d,'%s', '%s',STR_TO_DATE('%s', '%s'), '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d,'%s','%s', %d,%d, '%s')", u.RestaurantID, u.Name, u.Surname, u.Birthdate, date, u.Sex, hash, u.Email, u.Address, u.City, u.Country, u.PostalCode, u.Confirmed, u.VerifyCode, u.Role, u.ManagerID, u.AlwaysLogged, time.Now())
		_, err = db.Exec(statement)
		if err != nil {
			return err
		}
		err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)
		if err != nil {
			return err
		}
		verhash := randSeq(45)
		//log.Println(verhash)
		statement = fmt.Sprintf("UPDATE users SET verifycode = '%s' WHERE id=%d", verhash, u.ID)
		_, err = db.Exec(statement)
		if err != nil {
			return err
		}
		statement = fmt.Sprintf("SELECT password, verifycode, reg_date FROM users WHERE id=%d", u.ID)
		log.Println(statement)
		err = db.QueryRow(statement).Scan(&u.Password, &u.VerifyCode, &u.RegDate)
		if err != nil {
			return err
		}
		log.Println("URL: http://" + config.Owner.URL + "/user/" + strconv.Itoa(u.ID) + "/" + u.VerifyCode + "")
		mail.Send(u.Email, "Automated email from syncyours", "<strong>test: </strong><a href='http://"+config.Owner.URL+"/verify/"+strconv.Itoa(u.ID)+"/"+u.VerifyCode+"'>Verify Account</a>")

	default:

		err = errors.New("User Already exists")
		//log.Println("hi", err)
		return err

	}

	return nil
}
func getUsers(db *sql.DB, start, count int) ([]users, error) {
	statement := fmt.Sprintf("SELECT id, name, birthdate FROM users LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	user := []users{}
	for rows.Next() {
		var u users
		if err := rows.Scan(&u.ID, &u.Name, &u.Birthdate); err != nil {
			return nil, err
		}
		user = append(user, u)
	}
	return user, nil
}
