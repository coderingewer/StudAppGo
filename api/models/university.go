package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type University struct {
	gorm.Model
	Name      string    `json:"name"`
	CityID    uint      `json:"cityId"`
	UniMail   string    `json:"unimail"`
	Location  City      `json:"location"`
	Faculties []Faculty `gorm:"many2many:university_facultites" json:"faculties"`
}

type UniverstyFaculty struct {
	UniversityID uint    `gorm:"primary_key column:university_id"json:"universityId"`
	FacultyID    uint    `json:"facultyId"`
	Ffaculty     Faculty `json:"faculty"`
}

func (uni *University) Prepare() {
	uni.ID = 0
	uni.Location = City{}
	uni.Name = html.EscapeString(strings.TrimSpace(uni.Name))
}

func (uni *University) SAve() (*University, error) {
	err := GetDB().Debug().Create(&uni).Error
	if err != nil {
		return &University{}, err
	}
	if uni.ID != 0 {
		err := GetDB().Debug().Table("cities").Where("id=?", uni.CityID).Take(&uni.Location).Error
		if err != nil {
			return &University{}, err
		}
	}
	return uni, nil
}

func (uni *University) FindByCityID(cid uint) ([]University, error) {
	universities := []University{}
	db := GetDB().Debug().Table("universities").Where("city_id=?", cid).Limit(100).Find(&universities)
	if db.Error != nil {
		return []University{}, db.Error
	}
	if len(universities) > 0 {
		for i, _ := range universities {
			err := db.Debug().Table("cities").Where("id = ?", universities[i].CityID).Take(&universities[i].Location).Error
			if err != nil {
				return []University{}, err
			}
		}
	}
	return universities, nil
}

func (uni *University) FindBYID(unid uint) (*University, error) {
	err := GetDB().Debug().Table("universities").Where("id=?", unid).Take(&uni).Error
	if err != nil {
		return &University{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &University{}, errors.New("Universite bulunamadı")
	}
	err = db.Debug().Table("cities").Where("id = ?", uni.CityID).Take(&uni.Location).Error
	if err != nil {
		return &University{}, err
	}
	return uni, nil
}

func (unif UniverstyFaculty) AddAFacultyByID(unid, fid uint) (*UniverstyFaculty, error) {
	facultyUni := FacultyUniversity{}
	facultyUni.UniversityID = unid
	facultyUni.FacultyID = fid
	unif.UniversityID = unid
	unif.FacultyID = fid

	err := GetDB().Debug().Create(&unif).Error
	if err != nil {
		return &UniverstyFaculty{}, err
	}
	err = GetDB().Debug().Create(&facultyUni).Error
	if err != nil {
		return &UniverstyFaculty{}, err
	}
	if unif.UniversityID != 0 {
		db := GetDB().Debug().Table("faculties").Where("id=?", unif.FacultyID).Take(unif.Ffaculty)
		if db.Error != nil {
			return &UniverstyFaculty{}, db.Error
		}
	}
	if facultyUni.FacultyID != 0 {
		db := GetDB().Debug().Table("universities").Where("id=?", facultyUni.FacultyID).Take(facultyUni.Uuniversity)
		if db.Error != nil {
			return &UniverstyFaculty{}, db.Error
		}
	}
	return &unif, nil
}

func (unif UniverstyFaculty) GetFacultyByUniID(unid uint) ([]UniverstyFaculty, error) {
	faculties := []UniverstyFaculty{}
	err := GetDB().Debug().Table("university_facultites").Where("university_id= ?", unid).Find(&faculties).Error
	if err != nil {
		return []UniverstyFaculty{}, err
	}
	if len(faculties) > 0 {
		for i, _ := range faculties {
			err := db.Debug().Table("faculties").Where("id = ?", faculties[i].FacultyID).Take(&faculties[i].Ffaculty).Error
			if err != nil {
				return []UniverstyFaculty{}, err
			}
		}
	}
	return faculties, nil
}
