package models

import (
	"crypto/sha256"
	"encoding/json"
	"encoding/hex"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"regexp"
	"strings"

	db "github.com/ahmadkaisardev/databases"
	handler "github.com/ahmadkaisardev/handlers"
	middleware "github.com/ahmadkaisardev/middlewares"
)

type Auth struct {}

type AuthResp struct {
	Info string `json:"info"`
	Token string `json:"token"`
}

type JWTResp struct {
	Info string `json:"info"`
	Jwt jwt.MapClaims `json:"jwt"`
}

func (auth Auth) Response(w http.ResponseWriter, r *http.Request, status int, info string, token string) {
	var response AuthResp
	
	response.Info = info
	response.Token = token
	
	handler.Header(w)

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func (auth Auth) JResponse(w http.ResponseWriter, r *http.Request, status int, info string, Jwt jwt.MapClaims) {
	var response JWTResp
	
	response.Info = info
	response.Jwt = Jwt

	handler.Header(w)

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}


func (auth Auth) Authentication(w http.ResponseWriter, r *http.Request) {
	var jwt middleware.JWT

	var user User
	var users []User

	if r.Method == "OPTIONS" {
		handler.Response(w, r, 200, "allowed")
		return
	}

	email := r.FormValue("email")
	match, err := regexp.MatchString(`(^[a-zA-Z0-9\.\_\-]+)@([a-zA-Z0-9\.\_\-]+$)`, email)
	if !match || err != nil {
		handler.Response(w, r, 500, "invalid email")
		log.Println(err)
		return
	}

	sum := sha256.New()
	sum.Write([]byte(r.FormValue("password")))
	password_hash := hex.EncodeToString(sum.Sum(nil))
	password := password_hash

	rows, err := db.Gorm.Table("user").Select("user.id, user.name, role.id, role.name").Joins("JOIN role ON user.role_id = role.id").Where("user.email = ? AND user.password = ?", email, password).Limit(1).Rows()
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Role.Id, &user.Role.Name); err != nil{
			handler.Response(w, r, 500, "scan error")
			log.Println(err.Error())
			return
		} else{
			users = append(users, user)
		}
	}

	if len(users) < 1 {
		handler.Response(w, r, 500, "failed")
		log.Println("user not found")
		return
	}
	
	token, err := jwt.Generate(users[0].Id, users[0].Role.Id, users[0].Name, users[0].Role.Name)
	if err != nil {
		handler.Response(w, r, 500, "failed")
		log.Println("failed generate jwt")
		return
	}

	auth.Response(w, r, 200, "success", token)
}

func (auth Auth) Get(r *http.Request) (string, error) {

	var authorization string
	var bearer []string

	if authorization = r.Header.Get("Authorization"); len(authorization) > 0 {
		if bearer = strings.Split(authorization, " "); len(bearer) != 2 {
			return "", fmt.Errorf("invalid authorization")
		}
		return bearer[1], nil
	}
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		return "", fmt.Errorf("not authorized")
	} else {
		if bearer = strings.Split(cookie.Value, " "); len(bearer) != 2 {
			return "", fmt.Errorf("invalid authorization")
		}
		return bearer[1], nil
	}
	return "", fmt.Errorf("not authorized")
}

func (auth Auth) Claim(w http.ResponseWriter, r *http.Request) {
	var jwt middleware.JWT

	if r.Method == "OPTIONS" {
		handler.Response(w, r, 200, "allowed")
		return
	}

	token, err := auth.Get(r)
	if err != nil {
		handler.Response(w, r, 401, "not authorized")
		return
	}

	claims, err := jwt.Claim(token)
	if err != nil {
		handler.Response(w, r, 500, "error on claim")
		log.Println(err)
		return
	}

	auth.JResponse(w, r, 200, "success", claims)
}

func (auth Auth) Renew(w http.ResponseWriter, r *http.Request) {
	var jwt middleware.JWT
	var user User
	var users []User

	if r.Method == "OPTIONS" {
		handler.Response(w, r, 200, "allowed")
		return
	}

	token, err := auth.Get(r)
	if err != nil {
		handler.Response(w, r, 401, "not authorized")
		return
	}

	claims, err := jwt.Claim(token)
	if err != nil {
		handler.Response(w, r, 500, "error on claim")
		log.Println(err)
		return
	}

	rows, err := db.Gorm.Table("user").Select("user.id, user.name, role.id, role.name").Joins("JOIN role ON user.role_id = role.id").Where("user.id = ?", int(claims["user_id"].(float64))).Limit(1).Rows()
	if err != nil {
		handler.Response(w, r, 500, "query error")
		log.Println(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Name, &user.Role.Id, &user.Role.Name); err != nil{
			handler.Response(w, r, 500, "scan error")
			log.Println(err.Error())
			return
		} else{
			users = append(users, user)
		}
	}

	if len(users) < 1 {
		handler.Response(w, r, 500, "failed")
		log.Println("user not found")
		return
	}

	new_token, err := jwt.Generate(users[0].Id, users[0].Role.Id, users[0].Name, users[0].Role.Name)
	if err != nil {
		handler.Response(w, r, 500, "error renew token")
		log.Println(err)
		return
	}

	auth.Response(w, r, 200, "success", new_token);
}
