package main

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/context"
)

type JwtToken struct {
	UserID       int    `json:"user_id"`
	RestaurantID int    `json:"restaurant_id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Role         string `json:"role"`
	Token        string `json:"token"`
	Exp          int64  `json:"exp"`
}

type UserInfo struct {
	ID   int
	Name string
	Role string
}

type UserClaims struct {
	*jwt.StandardClaims
	TokenType string
	UserInfo
}

const (
	privKeyPath = "./keys/sample_key"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "./keys/sample_key.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func (j *JwtToken) Initialize() {
	var err error

	/////////////////////initialize signing keys//////////////////////////////
	signBytes, err := ioutil.ReadFile(config.JWT.PrivKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(config.JWT.PubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
	/////////////////////////////////////////////////////////////////////////
}

func createToken(user users) JwtToken {

	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	var timer int64
	if user.AlwaysLogged == 1 {
		timer = time.Now().Add(time.Minute * time.Duration(config.JWT.MaxTime)).Unix()
		log.Println("logging in for 24 hours", timer)
	} else {
		timer = time.Now().Add(time.Minute * time.Duration(config.JWT.DefaultTime)).Unix()
		log.Println("logging in for 2 hours", timer)
	}

	// set our claims
	t.Claims = &UserClaims{
		&jwt.StandardClaims{
			// set the expire time
			// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
			ExpiresAt: timer,
		},
		"RS256",
		UserInfo{user.ID, user.Name, user.Role},
	}

	// Creat token string
	token, _ := t.SignedString(signKey)
	/* if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Sorry, error while Signing Token!")
		return
	} */
	return JwtToken{UserID: user.ID, RestaurantID: user.RestaurantID, Email: user.Email, Name: user.Name, Surname: user.Surname, Role: user.Role, Token: token, Exp: timer}
	//return t.SignedString(signKey)
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	log.Println("Middleware called.................")
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		log.Println("header:", authorizationHeader)
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			log.Println("Bearer token length:", len(bearerToken))
			if len(bearerToken) == 2 {
				/* token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					log.Println("Bearer token:", bearerToken[1])
					 if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					//return []byte("secret"), nil
					return verifyKey, nil
				})
				if error != nil {
					//json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					respondWithError(w, http.StatusNotFound, error.Error())
					return
				} */
				token, err := request.ParseFromRequestWithClaims(req, request.OAuth2Extractor, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
					// since we only use the one private key to sign the tokens,
					// we also only use its public counter part to verify
					return verifyKey, nil
				})

				// If the token is missing or invalid, return error
				if err != nil {
					respondWithError(w, http.StatusUnauthorized, "Invalid authorization token")

				} else {
					if token.Valid {
						//log.Println("claim id :", token.Claims)
						log.Println("claim id :", token.Claims.(*UserClaims).ID)
						context.Set(req, "decoded", token.Claims.(*UserClaims).ID)
						next(w, req)
					} else {
						//json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
						respondWithError(w, http.StatusNotFound, "Invalid authorization token")
					}
				}
			}
		} else {
			//json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
			respondWithError(w, http.StatusNotFound, "An authorization header is required")
		}
	})
}
