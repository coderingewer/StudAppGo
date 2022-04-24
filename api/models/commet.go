package models

import (
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	UserID uint   `gorm:"not null" json:"userId"`
	PostID uint   `gorm:"not null" json:"postId"`
	Text   string `gorm:"not null;" json:"text"`
	User   User   `json:"user"`
}

func (cmmt *Comment) Prepare() {
	cmmt.UserID = 0
	cmmt.PostID = 0
	cmmt.User = User{}
	cmmt.Text = html.EscapeString(strings.TrimSpace(cmmt.Text))
}

func (cmmt *Comment) CreateCommet() (*Comment, error) {
	err := db.Debug().Create(&cmmt).Error
	if err != nil {
		return &Comment{}, err
	}
	if cmmt.PostID != 0 {
		err = db.Debug().Table("users").Where("user_id", cmmt.UserID).Take(cmmt.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return cmmt, nil
}

/*
func (cmmt *Comment) UpdateComment(pid uint) (*Comment, error) {
	db := db.Debug().Table("comments").Where("post_id = ?", pid).UpdateColumns(
		map[string]interface{}{
			"text": cmmt.Text,
		},
	)
	if db.Error != nil {
		return &Comment{}, db.Error
	}
	return cmmt, nil
}*/

func (cmmt *Comment) DeleteCommetByPostID(pid uint) (int64, error) {
	db := GetDB().Debug().Table("comments").Where("post_id = ? ", pid).Take(cmmt).Delete(Comment{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (cmmt *Comment) FindByUserID(uid uint) (*[]Comment, error) {
	comments := []Comment{}
	err := db.Debug().Table("comments").Where("user_id = ? ", uid).Find(comments).Error
	if err != nil {
		return &[]Comment{}, err
	}
	if len(comments) > 0 {
		for i, _ := range comments {
			err = db.Debug().Table("users").Where("id = ? ", comments[i].UserID).Find(comments[i].User).Error
			if err != nil {
				return &[]Comment{}, nil
			}
		}
	}
	return &comments, nil
}

func (cmmt *Comment) FindByPostID(pid uint) (*[]Comment, error) {
	comments := []Comment{}
	err := db.Debug().Table("comments").Where("post_id = ? ", pid).Find(comments).Error
	if err != nil {
		return &[]Comment{}, err
	}
	if len(comments) > 0 {
		for i, _ := range comments {
			err = db.Debug().Table("users").Where("id = ? ", comments[i].UserID).Find(comments[i].User).Error
			if err != nil {
				return &[]Comment{}, nil
			}
		}
	}
	return &comments, nil
}
