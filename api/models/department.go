package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Department struct {
	gorm.Model
	Name      string       `json:"name"`
	Faculties []Faculty    `gorm:"many2many:university_departments;" json:"-"`
	Schools   []University `gorm:"many2many:university_departments;" json:"-"`
}

type UniversityDepartment struct {
	UniversityID uint       `gorm:"primary_key column:university_id" json:"universityId"`
	FacultyID    uint       `gorm:"primary_key column:faculty_id" json:"facultyId"`
	DepartmentID uint       `gorm:"primary_key column:department_id" json:"departmentId"`
	University   University `json:"university"`
	Faculty      Faculty    `json:"faculty"`
	Department   Department `json:"department"`
}

func (d *Department) Prepare() {
	d.ID = 0
	d.Name = html.EscapeString(strings.TrimSpace(d.Name))
	d.Schools = []University{}
	d.Faculties = []Faculty{}
}

func (d *Department) Save() (*Department, error) {
	err := GetDB().Debug().Create(&d).Error
	if err != nil {
		return &Department{}, err
	}
	return d, nil
}

func (d *Department) FindAllDepartments() ([]Department, error) {
	departments := []Department{}
	db := GetDB().Table("departments").Limit(100).Find(&departments)
	if db.Error != nil {
		return []Department{}, db.Error
	}
	return departments, nil
}

func (d *Department) FindDepartmentByID(did uint) (*Department, error) {
	err := GetDB().Debug().Table("departments").Where("id=?", did).Take(&d).Error
	if err != nil {
		return &Department{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Department{}, errors.New("Department not found")
	}
	return d, nil
}

func (d *Department) DeleteByID(did uint) (int64, error) {
	db := GetDB().Debug().Table("faculties").Where("id=? ", did).Take(&d).Delete(&Department{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}

func (duni *UniversityDepartment) DeleteByID(dunid uint) (int64, error) {
	db := GetDB().Debug().Table("faculties").Where("id=? ", dunid).Take(&duni).Delete(&UniversityDepartment{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}

func (duni *UniversityDepartment) FindDepartmentByFacultyIDAndUniID(unid, fid uint) ([]UniversityDepartment, error) {
	uniDepartments := []UniversityDepartment{}
	db := GetDB().Table("university_departments").Where("faculty_id = ? AND university_id =?", duni.FacultyID, duni.UniversityID).Limit(100).Find(&uniDepartments)
	if db.Error != nil {
		return []UniversityDepartment{}, db.Error
	}
	if len(uniDepartments) > 0 {
		for i, _ := range uniDepartments {
			err := GetDB().Debug().Table("faculties").Where("id=?", uniDepartments[i].FacultyID).Take(uniDepartments[i].Faculty).Error
			if err != nil {
				return []UniversityDepartment{}, err
			}
			err = GetDB().Debug().Table("universities").Where("id=?", uniDepartments[i].UniversityID).Take(uniDepartments[i].University).Error
			if err != nil {
				return []UniversityDepartment{}, err
			}
			err = GetDB().Debug().Table("departments").Where("id=?", uniDepartments[i].DepartmentID).Take(uniDepartments[i].Department).Error
			if err != nil {
				return []UniversityDepartment{}, err
			}
		}
	}
	return uniDepartments, nil
}
