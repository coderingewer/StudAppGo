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
	UserID      uint      `gorm:"not null" json:"userId"`
	User        User      `json:"user"`
	CityID      uint      `gorm:"not null" json:"cityId"`
	City        City      `json:"city"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
}

func (amg Amigo) Prepare() {
	amg.Title = html.EscapeString(strings.TrimSpace(amg.Title))
	amg.Description = html.EscapeString(strings.TrimSpace(amg.Description))
	amg.UserID = 0
	amg.User = User{}
	amg.Deadline = time.Time{}
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

func (amg *Amigo) DeleteAmigoByUserID(aid uint) (int64, error) {
	db := db.Debug().Table("amigos").Where("user_id = ? ", aid).Take(&amg).Delete(&Amigo{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, db.Error
}

func (amg *Amigo) FindAllAmigosByDESC() (*[]Amigo, error) {
	amigos := []Amigo{}
	err := db.Table("amigos").Order("created_at desc").Find(amigos).Error
	if err != nil {
		return &[]Amigo{}, err
	}
	if len(amigos) > 0 {
		for i, _ := range amigos {
			err := db.Table("users").Where("id = ?", amigos[i].UserID).Take(amigos[i].User).Error
			if err != nil {
				return &[]Amigo{}, err
			}
			err = db.Table("cities").Where("id = ?", amigos[i].CityID).Take(amigos[i].City).Error
			if err != nil {
				return &[]Amigo{}, err
			}
		}
	}
	return &amigos, nil
}

func (amg *Amigo) FindByUserID(uid uint) (*[]Amigo, error) {
	amigos := []Amigo{}
	err := db.Table("amigos").Where("user_id = ?", uid).Find(amigos).Error
	if err != nil {
		return &[]Amigo{}, err
	}
	if len(amigos) > 0 {
		for i, _ := range amigos {
			err := db.Table("users").Where("id = ?", amigos[i].UserID).Take(amigos[i].User).Error
			if err != nil {
				return &[]Amigo{}, err
			}
			err = db.Table("cities").Where("id = ?", amigos[i].CityID).Take(amigos[i].City).Error
			if err != nil {
				return &[]Amigo{}, err
			}
		}
	}
	return &amigos, nil
}

func (amg *Amigo) FindByCityID(cid uint) (*[]Amigo, error) {
	amigos := []Amigo{}
	err := db.Table("amigos").Where("city_id = ?", cid).Find(amigos).Error
	if err != nil {
		return &[]Amigo{}, err
	}
	if len(amigos) > 0 {
		for i, _ := range amigos {
			err := db.Table("users").Where("id = ?", amigos[i].UserID).Take(amigos[i].User).Error
			if err != nil {
				return &[]Amigo{}, err
			}
			err = db.Table("cities").Where("id = ?", amigos[i].CityID).Take(amigos[i].City).Error
			if err != nil {
				return &[]Amigo{}, err
			}
		}
	}
	return &amigos, nil
}
