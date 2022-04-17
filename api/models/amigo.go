package models

import (
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Amigo struct {
	gorm.Model
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"not null" json:"descriptiom"`
	UserID      uint      `gorm:"primary_key; not null" json:"userId"`
	User        User      `json:"user"`
	Deatline    time.Time `gorm:"not null" json:"deatline"`
}

func (amg Amigo) Prepare() {
	amg.Title = html.EscapeString(strings.TrimSpace(amg.Title))
	amg.Description = html.EscapeString(strings.TrimSpace(amg.Description))
	amg.UserID = 0
	amg.User = User{}
	amg.Deatline = time.Time{}
}

func (amg *Amigo) CreateAmigo() (*Amigo, error) {
	err := db.Debug().Create(amg).Error
	if err != nil {
		return &Amigo{}, err
	}
	if amg.UserID != 0 {
		err := db.Debug().Table("users").Where("id = ? ", amg.UserID).Take(amg.User).Error
		if err != nil {
			return &Amigo{}, err
		}
	}
	return amg, nil
}

func (amg *Amigo) DeleteAmigoByID(aid uint) (int64, error) {
	db := db.Debug().Table("amigos").Where("id = ? ", aid).Take(&amg).Delete(&Amigo{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, db.Error
}
