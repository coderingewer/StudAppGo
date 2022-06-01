package models

import (
	"errors"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username     string     `gorm:"size:255;not null; unique" json:"username"`
	Email        string     `gorm:"size:255;not null;unique" json:"email"`
	Firstname    string     `json:"firstname"`
	Lastname     string     `json:"lastname"`
	UserRole     string     `gorm:"size:20;not null;" json:"userRole"`
	Password     string     `gorm:"size:255;not null;" json:"password"`
	UniversityID uint       `json:"-"`
	University   University `json:"-"`
	FacultyID    uint       `json:"-"`
	Faculty      Faculty    `json:"-"`
	DepartmentID uint       `gorm:"not null" json:"-"`
	Department   Department `json:"-"`
	Friends      []User     `gorm:"many2many:user_friends" json:"-"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashesPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashesPassword), []byte(password))
}

func (u *User) BeforeSAve() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		fmt.Println(string(err.Error()))
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Lastname = html.EscapeString(strings.TrimSpace(u.Lastname))
	u.Firstname = html.EscapeString(strings.TrimSpace(u.Lastname))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.UserRole = "USER"
	u.DeletedAt = nil
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Password == "" {
			return errors.New("Şifre Zorunlu")
		}
		if u.Email == "" {
			return errors.New("Email Zorunlu")
		}
		/*if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Doğrulanmamış E Mail Adresi")
		}*/
		return nil

	default:
		if u.Username == "" {
			return errors.New("Kullanıcı Adı Zorulu")
		}
		if u.Password == "" {
			return errors.New("Şifre Adı Zorulu")
		}
		if u.Username == "" {
			return errors.New("E Posta Adresi Zorulu")
		}
		if edu := strings.HasSuffix(u.Email, ".edu.tr"); !edu {
			return errors.New("Email üniversite Maili Olmalıdır")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Onaylanmamış E Mail Adresi")
		}
		return nil
	}
}

func (u *User) SaveUser() (*User, error) {
	db := GetDB()
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	/*err = GetDB().Debug().Table("universities").Where("id=?", u.UniversityID).Take(&u.University).Error
	if err != nil {
		return &User{}, err
	}
	err = GetDB().Debug().Table("faculties").Where("id=?", u.FacultyID).Take(&u.Faculty).Error
	if err != nil {
		return &User{}, err
	}*/
	return u, nil
}

func (u *User) FindAllUsers() ([]User, error) {
	var err error
	db := GetDB()
	users := []User{}
	err = db.Debug().Table("users").Limit(100).Find(&users).Error
	if err != nil {
		return []User{}, err
	}
	if len(users) > 0 {
		for i, _ := range users {
			err := GetDB().Debug().Table("universities").Where("id=?", users[i].UniversityID).Take(&users[i].University).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("faculties").Where("id=?", users[i].FacultyID).Take(&users[i].Faculty).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("cities").Where("id=?", users[i].University.CityID).Take(&users[i].University.Location).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("departments").Where("id=?", users[i].DepartmentID).Take(&users[i].Department).Error
			if err != nil {
				return []User{}, err
			}
		}
	}
	return users, nil
}

func (u *User) FindByUniversityID(unid uint) ([]User, error) {
	users := []User{}
	err := GetDB().Table("users").Where("university_id=?").Limit(100).Find(&users).Error
	if err != nil {
		return []User{}, err
	}
	if len(users) > 0 {
		for i, _ := range users {
			err := GetDB().Debug().Table("universities").Where("id=?", users[i].UniversityID).Take(&users[i].University).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("faculties").Where("id=?", users[i].FacultyID).Take(&users[i].Faculty).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("cities").Where("id=?", users[i].University.CityID).Take(&users[i].University.Location).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("departments").Where("id=?", users[i].DepartmentID).Take(&users[i].Department).Error
			if err != nil {
				return []User{}, err
			}
		}
	}
	return users, nil
}

func (u *User) FindByFacultyID(unid uint) ([]User, error) {
	users := []User{}
	err := GetDB().Table("users").Where("faculty_id=?").Limit(100).Find(&users).Error
	if err != nil {
		return []User{}, err
	}
	if len(users) > 0 {
		for i, _ := range users {
			err := GetDB().Debug().Table("universities").Where("id=?", users[i].UniversityID).Take(&users[i].University).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("faculties").Where("id=?", users[i].FacultyID).Take(&users[i].Faculty).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("cities").Where("id=?", users[i].University.CityID).Take(&users[i].University.Location).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("departments").Where("id=?", users[i].DepartmentID).Take(&users[i].Department).Error
			if err != nil {
				return []User{}, err
			}
		}
	}
	return users, nil
}

func (u *User) FindByDepartmentID(unid uint) ([]User, error) {
	users := []User{}
	err := GetDB().Table("users").Where("department_id=?").Limit(100).Find(&users).Error
	if err != nil {
		return []User{}, err
	}
	if len(users) > 0 {
		for i, _ := range users {
			err := GetDB().Debug().Table("universities").Where("id=?", users[i].UniversityID).Take(&users[i].University).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("faculties").Where("id=?", users[i].FacultyID).Take(&users[i].Faculty).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("cities").Where("id=?", users[i].University.CityID).Take(&users[i].University.Location).Error
			if err != nil {
				return []User{}, err
			}
			err = GetDB().Debug().Table("departments").Where("id=?", users[i].DepartmentID).Take(&users[i].Department).Error
			if err != nil {
				return []User{}, err
			}
		}
	}
	return users, nil
}

func (u *User) FindByID(uid uint) (*User, error) {
	var err error
	db = GetDB()
	err = db.Debug().Table("users").Where("id=?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("Kullanıcı Bulunamadı")
	}
	err = GetDB().Table("universities").Where("id =? ", u.UniversityID).Take(&u.University).Error
	if err != nil {
		return &User{}, err
	}
	err = GetDB().Debug().Table("faculties").Where("id=?", u.FacultyID).Take(&u.Faculty).Error
	if err != nil {
		return &User{}, err
	}
	err = GetDB().Debug().Table("departments").Where("id=?", u.DepartmentID).Take(&u.Department).Error
	if err != nil {
		return &User{}, err
	}
	return u, err
}

func (u *User) UpdateAUser(uid uint) (*User, error) {
	err := u.BeforeSAve()
	if err != nil {
		log.Fatal(err)
	}
	db := GetDB().Table("users").Where("id=?", uid).UpdateColumn(
		map[string]interface{}{
			"username":   u.Username,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	err = GetDB().Table("users").Where("id=?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteByID(uid uint) (int64, error) {
	db := GetDB().Debug().Table("users").Where("id=?", uid).Take(&u).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
