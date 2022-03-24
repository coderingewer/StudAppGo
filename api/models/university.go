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
	Location  City      `json:"location"`
	Faculties []Faculty `gorm:"many2many:university_faculties" json:"faculties"`
}

type UniverstyFaculty struct {
	UniversityID uint       `gorm:"primary_key column:university_id" json:"universityId"`
	FacultyID    uint       `json:"facultyId"`
	Faculty      Faculty    `json:"faculty"`
	University   University `json:"university"`
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
		return &University{}, errors.New("University not found")
	}
	err = db.Debug().Table("cities").Where("id = ?", uni.CityID).Take(&uni.Location).Error
	if err != nil {
		return &University{}, err
	}
	return uni, nil
}

func (unif UniverstyFaculty) AddAFacultyByID(unid, fid uint) (*UniverstyFaculty, error) {
	unif.UniversityID = unid
	unif.FacultyID = fid

	err := GetDB().Debug().Create(&unif).Error
	if err != nil {
		return &UniverstyFaculty{}, err
	}
	if unif.UniversityID != 0 {
		err = GetDB().Debug().Table("faculties").Where("id=?", unif.FacultyID).Take(unif.Faculty).Error
		if err != nil {
			return &UniverstyFaculty{}, err
		}
		err = GetDB().Debug().Table("universities").Where("id=?", unif.UniversityID).Take(unif.University).Error
		if err != nil {
			return &UniverstyFaculty{}, err
		}
	}
	return &unif, nil
}

func (duni UniversityDepartment) AddADepartmentByID(unid, did, fid uint) (*UniversityDepartment, error) {
	duni.UniversityID = unid
	duni.FacultyID = fid
	duni.DepartmentID = did

	err := GetDB().Debug().Create(&duni).Error
	if err != nil {
		return &UniversityDepartment{}, err
	}
	if duni.UniversityID != 0 {
		db := GetDB().Debug().Table("faculties").Where("id=?", duni.FacultyID).Take(duni.Faculty)
		if db.Error != nil {
			return &UniversityDepartment{}, db.Error
		}
		db = GetDB().Debug().Table("departments").Where("id=?", duni.DepartmentID).Take(duni.Department)
		if db.Error != nil {
			return &UniversityDepartment{}, db.Error
		}
		db = GetDB().Debug().Table("universities").Where("id=?", duni.UniversityID).Take(duni.University)
		if db.Error != nil {
			return &UniversityDepartment{}, db.Error
		}
	}
	return &duni, nil
}

func (uni *University) DeleteByID(unid uint) (int64, error) {
	db := GetDB().Debug().Table("faculties").Where("id=? ", unid).Take(&uni).Delete(&University{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}

func (unif *UniverstyFaculty) DeleteUniversityFacultyByID(unid uint) (int64, error) {
	db := GetDB().Debug().Table("university_facultites").Where("university_id=?", unid).Take(&unif).Delete(&UniverstyFaculty{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}

func (duni *UniversityDepartment) DeleteUniversityDepartmentByUniID(unid uint) (int64, error) {
	db := GetDB().Debug().Table("university_departments").Where("university_id=? ", unid).Take(&duni).Delete(&UniversityDepartment{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}
