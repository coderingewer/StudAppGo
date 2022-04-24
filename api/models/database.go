package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "uyumak"
	dbname   = "studapp"
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
		Image{}, UserFriend{}, FriendshipRequest{})

	fmt.Println("DB Successfully connected!")
}

func GetDB() *gorm.DB {
	return db
}
