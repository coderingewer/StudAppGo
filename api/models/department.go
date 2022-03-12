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
	Faculties []Faculty    `gorm:"many2many:department_universities;" json:"-"`
	Schools   []University `gorm:"many2many:department_universities;" json:"-"`
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

func (duni *UniversityDepartment) FindDepartmentByUniID() ([]UniversityDepartment, error) {
	departmentuni := []UniversityDepartment{}
	db := GetDB().Table("department_universities").Where("university_id", duni.UniversityID).Limit(100).Find(&departmentuni)
	if db.Error != nil {
		return []UniversityDepartment{}, db.Error
	}
	for i, _ := range departmentuni {
		err := GetDB().Debug().Table("faculties").Where("id=?", departmentuni[i].FacultyID).Take(departmentuni[i].Faculty).Error
		if err != nil {
			return []UniversityDepartment{}, err
		}
		err = GetDB().Debug().Table("universities").Where("id=?", departmentuni[i].UniversityID).Take(departmentuni[i].University).Error
		if err != nil {
			return []UniversityDepartment{}, err
		}
		err = GetDB().Debug().Table("departments").Where("id=?", departmentuni[i].DepartmentID).Take(departmentuni[i].Department).Error
		if err != nil {
			return []UniversityDepartment{}, err
		}
	}
	return departmentuni, nil
}

func (duni *UniversityDepartment) FindDepartmentByFacultyID() ([]UniversityDepartment, error) {
	departmentuni := []UniversityDepartment{}
	db := GetDB().Table("department_universities").Where("faculty_id", duni.FacultyID).Limit(100).Find(&departmentuni)
	if db.Error != nil {
		return []UniversityDepartment{}, db.Error
	}
	if len(departmentuni) > 0 {
		for i, _ := range departmentuni {
			err := GetDB().Debug().Table("faculties").Where("id=?", departmentuni[i].FacultyID).Take(departmentuni[i].Faculty).Error
			if err != nil {
				return []UniversityDepartment{}, err
			}
			err = GetDB().Debug().Table("universities").Where("id=?", departmentuni[i].UniversityID).Take(departmentuni[i].University).Error
			if err != nil {
				return []UniversityDepartment{}, err
			}
			err = GetDB().Debug().Table("departments").Where("id=?", departmentuni[i].DepartmentID).Take(departmentuni[i].Department).Error
			if err != nil {
				return []UniversityDepartment{}, err
			}
		}
	}
	return departmentuni, nil
}
