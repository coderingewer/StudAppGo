package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Faculty struct {
	gorm.Model
	Name string `json:"name"`
}

func (f *Faculty) Prepare() {
	f.ID = 0
	f.Name = html.EscapeString(strings.TrimSpace(f.Name))
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

func (faculty *Faculty) FindFacultyByID(fid uint) (*Faculty, error) {
	err := GetDB().Debug().Table("faculties").Where("id=?", fid).Take(&faculty).Error
	if err != nil {
		return &Faculty{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Faculty{}, errors.New("Faculty not found")
	}
	return faculty, nil
}

func (unif UniversityFaculty) GetFacultyByUniID(unid uint) ([]UniversityFaculty, error) {
	faculties := []UniversityFaculty{}
	err := GetDB().Debug().Table("university_facultites").Where("university_id= ?", unid).Find(&faculties).Error
	if err != nil {
		return []UniversityFaculty{}, err
	}
	if len(faculties) > 0 {
		for i, _ := range faculties {
			err := db.Debug().Table("faculties").Where("id = ?", faculties[i].FacultyID).Take(&faculties[i].Faculty).Error
			if err != nil {
				return []UniversityFaculty{}, err
			}
			err = db.Debug().Table("universities").Where("id = ?", faculties[i].UniversityID).Take(&faculties[i].University).Error
			if err != nil {
				return []UniversityFaculty{}, err
			}
		}
	}
	return faculties, nil
}

func (f *Faculty) DeleteByID(fid uint) (int64, error) {
	db := GetDB().Debug().Table("faculties").Where("id=? ", fid).Take(&f).Delete(&Faculty{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}
