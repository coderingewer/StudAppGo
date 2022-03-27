package models

import (
	"github.com/jinzhu/gorm"
)

type FriendshipRequest struct {
	gorm.Model
	SenderID   uint `gorm:"not null;" json:"senderId"`
	RecieverID uint `gorm:"not null;" json:"recieverId"`
	Sender     User `json:"sender"`
	Reciever   User `json:"reciever"`
}

func (fr *FriendshipRequest) Prepare() {
	fr.SenderID = 0
	fr.RecieverID = 0
	fr.Reciever = User{}
	fr.Sender = User{}
}

func (fr *FriendshipRequest) CreateRequest() (*FriendshipRequest, error) {

	err := GetDB().Debug().Create(&fr).Error
	if err != nil {
		return &FriendshipRequest{}, err
	}
	return fr, nil
}
