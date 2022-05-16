package databases

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var Gorm *gorm.DB

func GormConnect() (*gorm.DB, error) {
	MYSQL_HOST := os.Getenv("MYSQL_HOST")
	MYSQL_PORT := os.Getenv("MYSQL_PORT")
	MYSQL_USERNAME := os.Getenv("MYSQL_USERNAME")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_DATABASE := os.Getenv("MYSQL_DATABASE")
	Gorm, err := gorm.Open("mysql", MYSQL_USERNAME + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_HOST + ":"+ MYSQL_PORT + ")/" + MYSQL_DATABASE)
	if err != nil {
		log.Fatal(err)
	}
	Gorm.LogMode(true)

	return Gorm, err
}
