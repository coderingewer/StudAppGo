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

type UserFriend struct {
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

func (usrf *UserFriend) Prepare() {
	usrf.SenderID = 0
	usrf.RecieverID = 0
	usrf.Reciever = User{}
	usrf.Sender = User{}
}

func (fr *FriendshipRequest) CreateRequest() (*FriendshipRequest, error) {

	err := GetDB().Debug().Create(&fr).Error
	if err != nil {
		return &FriendshipRequest{}, err
	}
	return fr, nil
}

func (fr *FriendshipRequest) DeleteFrienshipRequestByUserID(frsid uint) (int64, error) {
	err := GetDB().Debug().Table("friendship_requests").Where("sender_id = ?", frsid).Take(&fr).Delete(&FriendshipRequest{}).Error
	if err != nil {
		return 0, err
	}
	return db.RowsAffected, nil
}

func (fr FriendshipRequest) GetRequestsByRecieverID(rid uint) (*[]FriendshipRequest, error) {
	requests := []FriendshipRequest{}
	err := db.Debug().Model(&FriendshipRequest{}).Where("reciever_id", rid).Find(&requests).Error
	if err != nil {
		return &[]FriendshipRequest{}, err
	}
	if len(requests) != 0 {
		for i, _ := range requests {
			err = db.Debug().Model(User{}).Where("id=?", requests[i].SenderID).Take(requests[i].Sender).Error
			if err != nil {
				return &[]FriendshipRequest{}, err
			}
		}
	}
	return &requests, nil
}

func (fr FriendshipRequest) GetRequestsBySenderID(rid uint) (*[]FriendshipRequest, error) {
	requests := []FriendshipRequest{}
	err := db.Debug().Model(&FriendshipRequest{}).Where("sender_id", rid).Find(&requests).Error
	if err != nil {
		return &[]FriendshipRequest{}, err
	}
	if len(requests) != 0 {
		for i, _ := range requests {
			err = db.Debug().Model(User{}).Where("id=?", requests[i].RecieverID).Take(requests[i].Reciever).Error
			if err != nil {
				return &[]FriendshipRequest{}, err
			}
		}
	}
	return &requests, nil
}

func (uf *UserFriend) CreateFriend() (*UserFriend, error) {
	err := db.Debug().Create(&uf).Error
	if err != nil {
		return &UserFriend{}, err
	}
	if uf.ID != 0 {
		err = db.Debug().Table("users").Where("id = ?", uf.SenderID).Take(uf.Sender).Error
		if err != nil {
			return &UserFriend{}, err
		}
	}
	return uf, nil
}

func (uf *UserFriend) FindFriendsByUserID(uid uint) (*[]UserFriend, error) {
	friends := []UserFriend{}
	err := db.Debug().Table("user_firends").Where("reciever_id", uid).Or(UserFriend{SenderID: uid}).Find(friends).Error
	if err != nil {
		return &[]UserFriend{}, nil
	}
	if len(friends) > 0 {
		for i, _ := range friends {
			if uid == friends[i].RecieverID {
				err = db.Debug().Table("users").Where("id=?", friends[i].SenderID).Take(friends[i].Sender).Error
				if err != nil {
					return &[]UserFriend{}, nil
				}
			}
			if uid == friends[i].SenderID {
				err = db.Debug().Table("users").Where("id=?", friends[i].RecieverID).Take(friends[i].Reciever).Error
				if err != nil {
					return &[]UserFriend{}, nil
				}
			}
		}
	}
	return &friends, nil
}
