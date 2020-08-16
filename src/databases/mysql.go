package databases

import (
	"database/sql"
  	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var SQL *sql.DB

func SQLConnect() (*sql.DB, error) {
	SQL, err := sql.Open("mysql", "<db_username>:<db_password>@tcp(localhost:3306)/<table_name>")
	if err != nil {
		log.Fatal(err)
	}

	return SQL, err
}
