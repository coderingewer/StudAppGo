package models

import (
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Sender   User   `json:"sender"`
	UserID   uint   `gorm:"not null" json:"userId"`
	PhotosID int    `json:"photosId"`
	Like     int    `json:"like"`
}

func (p *Post) Prepare() {
	p.ID = 0
	p.Like = 0
	p.Sender = User{}
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
}

func (p *Post) Save() (*Post, error) {
	err := db.Debug().Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?",
			p.UserID).Take(&p.Sender).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) FindAllPosts() ([]Post, error) {
	posts := []Post{}
	err := GetDB().Debug().Table("posts").Find(&posts).Error
	if err != nil {
		return []Post{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := GetDB().Debug().Table("users").Where("id=?", posts[i].UserID).Take(posts[i].Sender).Error
			if err != nil {
				return []Post{}, err
			}
		}
	}
	return posts, nil
}

func (p *Post) FindByID(pid uint) (*Post, error) {
	err := GetDB().Debug().Table("posts").Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Post{}, errors.New("Gönderi bulunamadı")
	}
	return p, nil
}

func (p *Post) UpdatePost(pid uint) (*Post, error) {
	db := GetDB().Debug().Table("posts").Where("id=?", pid).UpdateColumns(
		map[string]interface{}{
			"title": p.Title,
		},
	)
	if db.Error != nil {
		return &Post{}, db.Error
	}
	err := GetDB().Debug().Table("posts").Where("id=?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	return p, nil
}

func (p *Post) DeleteByID(pid uint) (int64, error) {
	db := GetDB().Debug().Table("posts").Where("id=? ", pid).Take(&p).Delete(&Post{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}

func (p *Post) FinBYUserID(uid uint) ([]Post, error) {
	posts := []Post{}
	err := GetDB().Debug().Table("posts").Where("user_id = ?", uid).Find(&posts).Error
	if err != nil {
		return []Post{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := GetDB().Debug().Table("users").Where("id=?", posts[i].UserID).Take(posts[i].Sender).Error
			if err != nil {
				return []Post{}, err
			}
		}
	}
	return posts, nil
}
