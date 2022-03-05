package models

import (
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Faculty struct {
	gorm.Model
	Name    string       `json:"name"`
	Schools []University `gorm:"many2many:faculty_universities;"json:"school"`
}

type FacultyUniversity struct {
	UniversityID uint       `gorm:"primary_key column:university_id"json:"universityId"`
	FacultyID    uint       `json:"facultyId"`
	Uuniversity  University `json:"faculty"`
}

func (f *Faculty) Prepare() {
	f.ID = 0
	f.Name = html.EscapeString(strings.TrimSpace(f.Name))
	f.Schools = []University{}
}

func (f *Faculty) Save() (*Faculty, error) {
	err := GetDB().Debug().Create(&f).Error
	if err != nil {
		return &Faculty{}, err
	}
	return f, nil
}

func (f *Faculty) FindAllFaculty() ([]Faculty, error) {
	faculties := []Faculty{}
	db := GetDB().Table("faculties").Limit(100).Find(&faculties)
	if db.Error != nil {
		return []Faculty{}, db.Error
	}
	return faculties, nil
}
