package models

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"regexp"
	"time"

	db "../databases"
	handler "../handlers"
	middleware "../middlewares"
)

type Cookie struct {}

func (cookie Cookie) Set(w http.ResponseWriter, name string, value string) {
	expire := time.Now().Add(time.Minute * 10)
	c := http.Cookie{
		Name: name,
		Value: value,
		Expires: expire,
		HttpOnly: true,
		Path:"/",
	}

	http.SetCookie(w, &c)
}

func (cookie Cookie) Authentication(w http.ResponseWriter, r *http.Request) {
	var err error
	var jwt middleware.JWT
	var token string
	var user User
	var users []User

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
		if err = rows.Scan(&user.Id, &user.Name, &user.Role.Id, &user.Role.Name); err != nil{
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
	
	token, err = jwt.Generate(users[0].Id, users[0].Role.Id, users[0].Name, users[0].Role.Name)
	if err != nil {
		handler.Response(w, r, 500, "failed")
		log.Println("failed generate jwt")
		return
	}

	cookie.Set(w, "Authorization", "Bearer " + token)
	handler.Response(w, r, 200, "success")
}

func (cookie Cookie) Renew(w http.ResponseWriter, r *http.Request) {
	var auth Auth
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

	cookie.Set(w, "Authorization", "Bearer " + new_token)
	handler.Response(w, r, 200, "success")
}

func (cookie Cookie) Destroy(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		handler.Response(w, r, 200, "allowed")
		return
	}

	c := http.Cookie{
		Name: "Authorization",
		Value: "Bearer ",
		Expires: time.Unix(0, 0),
		MaxAge: -1,
		HttpOnly: true,
		Path:"/",
	}

	http.SetCookie(w, &c)
	handler.Response(w, r, 200, "success")
}
