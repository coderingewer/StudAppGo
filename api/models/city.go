package models

import (
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type City struct {
	gorm.Model
	Name    string `json:"name"`
	Country string `json:"country"`
}

func (c *City) Prepare() {
	c.ID = 0
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.Country = html.EscapeString(strings.TrimSpace(c.Country))
}

func (c *City) Save() (*City, error) {
	err := GetDB().Debug().Create(&c).Error
	if err != nil {
		return &City{}, err
	}
	return c, nil
}

func (c *City) FindAll() ([]City, error) {
	cities := []City{}
	db := GetDB().Debug().Table("cities").Limit(100).Find(&cities)
	if db.Error != nil {
		return []City{}, db.Error
	}
	return cities, nil
}
