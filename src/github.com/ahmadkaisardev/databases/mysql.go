package databases

import (
	"database/sql"
  	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

var SQL *sql.DB

func SQLConnect() (*sql.DB, error) {
	MYSQL_HOST := os.Getenv("MYSQL_HOST")
	MYSQL_PORT := os.Getenv("MYSQL_PORT")
	MYSQL_USERNAME := os.Getenv("MYSQL_USERNAME")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_DATABASE := os.Getenv("MYSQL_DATABASE")
	SQL, err := sql.Open("mysql", MYSQL_USERNAME + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_HOST + ":" + MYSQL_PORT + ")/" + MYSQL_DATABASE)
	if err != nil {
		log.Fatal(err)
	}

	return SQL, err
}
