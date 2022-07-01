package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "ec2-52-49-120-150.eu-west-1.compute.amazonaws.com"
	port     = 5432
	user     = "yegnwdbtcidghq"
	password = "fb4de03fd8a6faffaa8a57684ee80cc903f96ba9cd9ee55d415ce735f0cd2ca5"
	dbname   = "d1j5rp6djesaba"
)

var db *gorm.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)
	conn, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	db = conn
	db.Debug().AutoMigrate(User{}, Post{}, Skill{},
		City{}, University{}, Faculty{},
		Department{}, UniversityDepartment{},
		UniversityFaculty{}, Comment{}, Amigo{},
		Image{}, Tag{})

	fmt.Println("DB Successfully connected!")
}

func GetDB() *gorm.DB {
	return db
}
