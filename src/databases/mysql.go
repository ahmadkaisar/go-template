package databases

import (
	"database/sql"
  	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var SQL *sql.DB

func SQLConnect() (*sql.DB, error) {
	SQL, err := sql.Open("mysql", "admin:M@R!4dB@tcp(localhost:3306)/go")
	if err != nil {
		log.Fatal(err)
	}

	return SQL, err
}
