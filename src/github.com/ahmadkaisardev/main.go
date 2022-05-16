package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	db "github.com/ahmadkaisardev/databases"
	handler "github.com/ahmadkaisardev/handlers"
	middleware "github.com/ahmadkaisardev/middlewares"
	model "github.com/ahmadkaisardev/models"
)

func main() {
	// connect to db using function connect ( see `mysql.go` )
	db.SQL, handler.Err = db.SQLConnect()
	db.Gorm, handler.Err = db.GormConnect()
	defer db.SQL.Close()
	defer db.Gorm.Close()

	// start http listener
	port := os.Getenv("PORT")
	router := mux.NewRouter()

	// middleware handler
	var jwt middleware.JWT

	// basic response handler
	router.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		handler.Response(w, r, 200, "success")
	}).Methods("GET", "OPTIONS")

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
	router.Handle("/users", jwt.Validate(user.Get, []int{1, 2})).Methods("GET", "OPTIONS")
	router.Handle("/users", jwt.Validate(user.Post, []int{1})).Methods("POST", "OPTIONS")
	router.Handle("/users/{id}", jwt.Validate(user.GetId, []int{1})).Methods("GET", "OPTIONS")
	router.Handle("/users/{id}", jwt.Validate(user.Put, []int{1})).Methods("PUT", "OPTIONS")
	router.Handle("/users/{id}", jwt.Validate(user.Delete, []int{1})).Methods("DELETE", "OPTIONS")

	// role handler
	var role model.Role
	router.Handle("/roles", jwt.Validate(role.Get, []int{1})).Methods("GET", "OPTIONS")
	router.Handle("/roles", jwt.Validate(role.Post, []int{})).Methods("POST", "OPTIONS")
	router.Handle("/roles/{id}", jwt.Validate(role.GetId, []int{1})).Methods("GET", "OPTIONS")
	router.Handle("/roles/{id}", jwt.Validate(role.Put, []int{})).Methods("PUT", "OPTIONS")
	router.Handle("/roles/{id}", jwt.Validate(role.Delete, []int{})).Methods("DELETE", "OPTIONS")

	// file handler
	var file model.File
	router.Handle("/files/file", jwt.Validate(file.Post, []int{1})).Methods("POST", "OPTIONS")
	router.Handle("/files/file/{id}", jwt.Validate(file.Put, []int{1})).Methods("PUT", "OPTIONS")
	router.Handle("/files/file/{id}", jwt.Validate(file.Delete, []int{1})).Methods("DELETE", "OPTIONS")
	router.Handle("/files/file/{name}", jwt.Validate(file.Get, []int{1})).Methods("GET", "OPTIONS")
	
	http.Handle("/", router)
	fmt.Println("Connected to port " + port)
	log.Fatal(http.ListenAndServe(":" + port, router))
}
