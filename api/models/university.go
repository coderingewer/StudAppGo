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

type UniversityFaculty struct {
	gorm.Model
	UniversityID uint       `gorm:"primary_key column:university_id" json:"universityId"`
	FacultyID    uint       `gorm:"primary_key column:faculty_id" json:"facultyId"`
	Faculty      Faculty    `json:"faculty"`
	University   University `json:"university"`
}

type UniversityDepartment struct {
	gorm.Model
	UniversityID uint       `gorm:"primary_key column:university_id" json:"universityId"`
	FacultyID    uint       `gorm:"primary_key column:faculty_id" json:"facultyId"`
	DepartmentID uint       `gorm:"primary_key column:department_id" json:"departmentId"`
	University   University `json:"university"`
	Faculty      Faculty    `json:"faculty"`
	Department   Department `json:"department"`
}

func (uni *University) Prepare() {
	uni.ID = 0
	uni.Location = City{}
	uni.Name = html.EscapeString(strings.TrimSpace(uni.Name))
}

func (unif *UniversityFaculty) Prepare() {
	unif.UniversityID = 0
	unif.FacultyID = 0
	unif.University = University{}
	unif.Faculty = Faculty{}
}

func (duni *UniversityDepartment) Prepare() {
	duni.UniversityID = 0
	duni.FacultyID = 0
	duni.University = University{}
	duni.Faculty = Faculty{}
	duni.DepartmentID = 0
	duni.Department = Department{}

}

func (uni *University) Save() (*University, error) {
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

func (unif *UniversityFaculty) AddAFacultyByID() (*UniversityFaculty, error) {

	err := db.Debug().Create(&unif).Error
	if err != nil {
		return &UniversityFaculty{}, err
	}
	if unif.ID != 0 {
		err := GetDB().Debug().Table("universities").Where("id=?", unif.UniversityID).Take(&unif.University).Error
		if err != nil {
			return &UniversityFaculty{}, err
		}
		if unif.University.ID != 0 {
			err := GetDB().Debug().Table("cities").Where("id=?", unif.University.CityID).Take(&unif.University.Location).Error
			if err != nil {
				return &UniversityFaculty{}, err
			}
		}
		err = GetDB().Debug().Table("faculties").Where("id=?", unif.FacultyID).Take(&unif.Faculty).Error
		if err != nil {
			return &UniversityFaculty{}, err
		}
	}
	return unif, nil
}

func (duni *UniversityDepartment) AddADepartmentByID() (*UniversityDepartment, error) {

	err := db.Debug().Create(&duni).Error
	if err != nil {
		return &UniversityDepartment{}, err
	}
	if duni.ID != 0 {
		err := GetDB().Debug().Table("universities").Where("id=?", duni.UniversityID).Take(&duni.University).Error
		if err != nil {
			return &UniversityDepartment{}, err
		}
		if duni.University.ID != 0 {
			err := GetDB().Debug().Table("cities").Where("id=?", duni.University.CityID).Take(&duni.University.Location).Error
			if err != nil {
				return &UniversityDepartment{}, err
			}
		}
		err = GetDB().Debug().Table("faculties").Where("id=?", duni.FacultyID).Take(&duni.Faculty).Error
		if err != nil {
			return &UniversityDepartment{}, err
		}

		err = GetDB().Debug().Table("departments").Where("id=?", duni.DepartmentID).Take(&duni.Department).Error
		if err != nil {
			return &UniversityDepartment{}, err
		}

	}
	return duni, nil
}

func (uni *University) DeleteByID(unid uint) (int64, error) {
	db := GetDB().Debug().Table("faculties").Where("id=? ", unid).Take(&uni).Delete(&University{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}

func (unif *UniversityFaculty) DeleteUniversityFacultyByID(unid uint) (int64, error) {
	db := GetDB().Debug().Table("university_facultites").Where("id=?", unid).Take(&unif).Delete(&UniversityFaculty{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}

func (duni *UniversityDepartment) DeleteUniversityDepartmentByUniID(unid uint) (int64, error) {
	db := GetDB().Debug().Table("university_departments").Where("id=? ", unid).Take(&duni).Delete(&UniversityDepartment{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}
