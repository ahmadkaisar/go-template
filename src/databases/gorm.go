package databases

import (
	"github.com/jinzhu/gorm"
	"log"
)

var Gorm *gorm.DB

func GormConnect() (*gorm.DB, error) {
	Gorm, err := gorm.Open("mysql", "<db_username>:<db_password>@tcp(localhost:3306)/<table_name>")
	if err != nil {
		log.Fatal(err)
	}
	Gorm.LogMode(true)

	return Gorm, err
}
