package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Skill struct {
	gorm.Model
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Scope    string `json:"scope"`
	AuthorID uint   `json:"author_id"`
	User     User   `json:"user"`
}

func (s *Skill) Prepare() {
	s.ID = 0
	s.Title = html.EscapeString(strings.TrimSpace(s.Title))
	s.Detail = html.EscapeString(strings.TrimSpace(s.Detail))
	s.Scope = html.EscapeString(strings.TrimSpace(s.Scope))
}

func (s *Skill) Save() (*Skill, error) {
	err := GetDB().Debug().Create(&s).Error
	if err != nil {
		return &Skill{}, err
	}
	return s, nil
}

func (s *Skill) FindAllSkills() ([]Skill, error) {
	skills := []Skill{}
	err := GetDB().Debug().Table("skills").Limit(100).Find(&skills).Error
	if err != nil {
		return []Skill{}, err
	}
	if len(skills) > 0 {
		for i, _ := range skills {
			err := GetDB().Debug().Table("users").Where("id=?", skills[i].AuthorID).Take(skills[i].User).Error
			if err != nil {
				return []Skill{}, err
			}
		}

	}
	return skills, nil
}

func (s *Skill) FindByID(sid uint) (*Skill, error) {
	err := GetDB().Debug().Table("skills").Where("id = ?", sid).Take(&s).Error
	if err != nil {
		return &Skill{}, err
	}
	if err == gorm.ErrRecordNotFound {
		return &Skill{}, errors.New("Bulunamadı")
	}
	err = GetDB().Debug().Table("users").Where("id=?", s.AuthorID).Take(s.User).Error
	if err != nil {
		return &Skill{}, err
	}
	return s, nil
}

func (s *Skill) UpdateSkill(sid uint) (*Skill, error) {
	db := GetDB().Debug().Table("skills").Where("id = ?", sid).UpdateColumns(
		map[string]interface{}{

			"title":  s.Title,
			"detail": s.Detail,
			"scope":  s.Scope,
		},
	)
	if db.Error != nil {
		return &Skill{}, db.Error
	}
	err := GetDB().Debug().Table("skills").Where("id = ?", sid).Take(&s).Error
	if err != nil {
		return &Skill{}, nil
	}
	err = GetDB().Debug().Table("users").Where("id=?", s.AuthorID).Take(s.User).Error
	if err != nil {
		return &Skill{}, err
	}
	return s, nil
}

func (s *Skill) DeleteByID(sid uint) (int64, error) {
	db := GetDB().Debug().Table("skills").Where("id=?", sid).Take(s).Delete(&Skill{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (s *Skill) FindByUserID(uid uint) ([]Skill, error) {
	skills := []Skill{}
	err := GetDB().Debug().Table("skills").Where("author_id=?", uid).Find(&skills).Error
	if err != nil {
		return []Skill{}, err
	}
	if len(skills) > 0 {
		for i, _ := range skills {
			err := GetDB().Debug().Table("users").Where("id=?", skills[i].AuthorID).Take(skills[i].User).Error
			if err != nil {
				return []Skill{}, err
			}
		}
	}

	return skills, nil
}
