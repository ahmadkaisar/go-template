package databases

import (
	"github.com/jinzhu/gorm"
	"log"
)

var Gorm *gorm.DB

func GormConnect() (*gorm.DB, error) {
	Gorm, err := gorm.Open("mysql", "admin:M@R!4dB@tcp(localhost:3306)/go")
	if err != nil {
		log.Fatal(err)
	}
	Gorm.LogMode(true)

	return Gorm, err
}
