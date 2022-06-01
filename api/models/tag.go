package models

type Tag struct {
	gorm.Model
	Tagname string `gorm:"not null" json:"tagname"`
	}

	
	func (tag *Tag) Prepare() {
		tag.UserID = 0
		tag.PostID = 0
		tag.User = User{}
		tag.Text = html.EscapeString(strings.TrimSpace(tag.Text))
	}
	
	func (tag *Tag) CreatTag() (*Tag, error) {
		err := db.Debug().Create(&tag).Error
		if err != nil {
			return &Tag{}, err
		}
		return tag, nil
	}
	