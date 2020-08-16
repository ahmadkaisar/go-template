package main

import (
	"fmt"
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	db "./databases"
	handler "./handlers"
	middleware "./middlewares"
	model "./models"
)

func main() {
	// connect to db using function connect ( see `mysql.go` )
	db.SQL, handler.Err = db.SQLConnect()
	db.Gorm, handler.Err = db.GormConnect()
	defer db.SQL.Close()
	defer db.Gorm.Close()

	// start http listener
	router := mux.NewRouter()

	// middleware handler
	var jwt middleware.JWT

	// auth
	var auth model.Auth
	router.HandleFunc("/auth", auth.Authentication).Methods("POST", "OPTIONS")
	router.HandleFunc("/auth/renew", auth.Renew).Methods("GET", "OPTIONS")
	router.HandleFunc("/claim", auth.Claim).Methods("GET", "OPTIONS")

	// cookie
	var cookie model.Cookie
	router.HandleFunc("/cookie/auth", cookie.Authentication).Methods("GET", "POST", "OPTIONS")
	router.HandleFunc("/cookie/auth/renew", cookie.Renew).Methods("GET", "OPTIONS")
	router.HandleFunc("/cookie/auth/destroy", cookie.Destroy).Methods("GET", "OPTIONS")

	// user handler
	var user model.User
	router.Handle("/users", jwt.Validate(user.Get)).Methods("GET", "OPTIONS")
	router.Handle("/users", jwt.Validate(user.Post)).Methods("POST", "OPTIONS")
	router.Handle("/users/{id}", jwt.Validate(user.GetId)).Methods("GET", "OPTIONS")
	router.Handle("/users/{id}", jwt.Validate(user.Put)).Methods("PUT", "OPTIONS")
	router.Handle("/users/{id}", jwt.Validate(user.Delete)).Methods("DELETE", "OPTIONS")

	// role handler
	var role model.Role
	router.Handle("/roles", jwt.Validate(role.Get)).Methods("GET", "OPTIONS")
	router.Handle("/roles", jwt.Validate(role.Post)).Methods("POST", "OPTIONS")
	router.Handle("/roles/{id}", jwt.Validate(role.GetId)).Methods("GET", "OPTIONS")
	router.Handle("/roles/{id}", jwt.Validate(role.Put)).Methods("PUT", "OPTIONS")
	router.Handle("/roles/{id}", jwt.Validate(role.Delete)).Methods("DELETE", "OPTIONS")

	// file handler
	var file model.File
	router.Handle("/files/file", jwt.Validate(file.Post)).Methods("POST", "OPTIONS")
	router.Handle("/files/file/{id}", jwt.Validate(file.Put)).Methods("PUT", "OPTIONS")
	router.Handle("/files/file/{id}", jwt.Validate(file.Delete)).Methods("DELETE", "OPTIONS")
	router.Handle("/files/file/{name}", jwt.Validate(file.Get)).Methods("GET", "OPTIONS")
	
	http.Handle("/", router)
	fmt.Println("Connected to port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}
