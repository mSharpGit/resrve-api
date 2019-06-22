package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, host, port, dbname string) {
	var err error

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.DB.SetConnMaxLifetime(1000 * time.Millisecond)
	//defer a.DB.Close()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(a.Router)))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
	a.Router.HandleFunc("/auth", a.authUsers).Methods("POST")
	a.Router.HandleFunc("/verify/{id:[0-9]+}/{code}", a.verifyUser).Methods("POST")
	a.Router.HandleFunc("/resetuser", a.resetUser).Methods("POST")
	a.Router.HandleFunc("/resetuser/{id:[0-9]+}/{code}", a.updateUserPass).Methods("POST")
	a.Router.HandleFunc("/reservations/{id:[0-9]+}", ValidateMiddleware(a.getReservations)).Methods("GET")
	a.Router.HandleFunc("/reservation", ValidateMiddleware(a.addReservation)).Methods("POST")
	a.Router.HandleFunc("/reservation/{id:[0-9]+}", ValidateMiddleware(a.deleteReservation)).Methods("DELETE")
	a.Router.HandleFunc("/section/floor/{id:[0-9]+}", ValidateMiddleware(a.getSections)).Methods("GET")
	a.Router.HandleFunc("/floor/restaurant/{id:[0-9]+}", ValidateMiddleware(a.getFloors)).Methods("GET")
	a.Router.HandleFunc("/table/section/{id:[0-9]+}", ValidateMiddleware(a.getTables)).Methods("GET")
	a.Router.HandleFunc("/user/restaurant/{id:[0-9]+}", ValidateMiddleware(a.getWaiters)).Methods("GET")
}

//////////////////////////////Tables
func (a *App) getWaiters(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Floor ID")
		return
	}
	waiter, err := getWaiters(a.DB, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, waiter)
}

//////////////////////////////Tables
func (a *App) getTables(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Floor ID")
		return
	}
	table, err := getTables(a.DB, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, table)
}

//////////////////////////////floors
func (a *App) getFloors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Floor ID")
		return
	}
	floor, err := getFloors(a.DB, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, floor)
}

//////////////////////////////sections
func (a *App) getSections(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Section ID")
		return
	}
	section, err := getSections(a.DB, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, section)
}

//////////////////////////////reservations
func (a *App) getReservations(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Restaurant ID")
		return
	}
	ocuppy, err := getReservations(a.DB, id, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, ocuppy)
}
func (a *App) addReservation(w http.ResponseWriter, r *http.Request) {
	var o occupy
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&o); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := o.addReservation(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, o)
}
func (a *App) deleteReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	o := occupy{ID: id}
	if err := o.deleteReservation(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

//////////////////////////////Users
func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	u := users{ID: id}
	if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	users, err := getUsers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var u users
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := u.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) verifyUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	code := vars["code"]
	var u users
	u.ID = id
	u.VerifyCode = code
	if err := u.verifyUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:

			respondWithError(w, http.StatusNotFound, "Wrong Validation Request")
		default:

			respondWithError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)

}
func (a *App) resetUser(w http.ResponseWriter, r *http.Request) {

	var u users
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := u.resetUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User does not exist")
		default:
			respondWithError(w, http.StatusNotFound, "Something went wrong with this request")
		}
		return
	}
	var reset resetpass
	reset.UserID = u.ID
	reset.Confirmed = 0
	reset.ResetDate = time.Now().String()
	reset.Code = randSeq(45)
	reset.OldPass = u.Password
	err := reset.addResetPass(a.DB)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Something went wrong with this request")
	}
	mail.Send(u.Email, "Automated email from syncyours To Reset Password", "<strong>Reset Password: </strong><a href='http://"+config.Owner.URL+"/resetuser/"+strconv.Itoa(reset.ID)+"/"+reset.Code+"'>RESET</a>")

	respondWithJSON(w, http.StatusCreated, u)
}
func (a *App) updateUserPass(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	code := vars["code"]
	reset := resetpass{ID: id, Code: code}

	//fetching the reset request
	if err := reset.getResetPass(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Reset Request not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reset); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	//log.Println("Pasword: ", reset.NewPass)
	hash, err := bcrypt.GenerateFromPassword([]byte(reset.NewPass), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid bcrypt failed")
	}

	reset.NewPass = string(hash)
	reset.Confirmed = 1
	if err := reset.updateNewPass(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:

			respondWithError(w, http.StatusNotFound, "New Password Update Request has been expired")
		default:

			respondWithError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}
	var u users
	u.ID = reset.UserID
	u.Password = string(hash)
	if err := u.updateUserPass(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:

			respondWithError(w, http.StatusNotFound, "User does not exist")
		default:

			respondWithError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}

	respondWithJSON(w, http.StatusOK, "UPDATED")

}
func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var u users
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	u.ID = id
	if err := u.updateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	u := users{ID: id}
	if err := u.deleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) authUsers(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var u users
	err = json.Unmarshal(b, &u)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := u.authUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:

			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			var ErrHashTooShort = errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password")
			if err.Error() == ErrHashTooShort.Error() {
				respondWithError(w, http.StatusBadRequest, "Password is wrong")
			} else {
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}

		}
		return
	}

	///adding a history record
	var l loginhist
	err = json.Unmarshal(b, &l)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload22")
		return
	}
	l.UserID = u.ID
	l.Flag = 1
	if err := l.addHist(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	//create the token
	tokenString := createToken(u)
	/* if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Sorry, error while Signing Token!")
		return
	} */

	respondWithJSON(w, http.StatusCreated, tokenString)
	//respondWithJSON(w, http.StatusCreated, u)

	/* if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	} */

}

// only accessible with a valid token
func restrictedHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from request
	token, err := request.ParseFromRequestWithClaims(r, request.OAuth2Extractor, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})

	// If the token is missing or invalid, return error
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token:"+err.Error())
		//w.WriteHeader(http.StatusUnauthorized)
		//fmt.Fprintln(w, "Invalid token:", err)
		return
	}

	// Token is valid
	respondWithJSON(w, http.StatusCreated, map[string]string{"name": token.Claims.(*UserClaims).Name})
	//fmt.Fprintln(w, "Welcome,", token.Claims.(*UserClaims).Name)
	return
}
