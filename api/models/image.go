package models

import (
	"html"
	"mime/multipart"
	"strings"
	"studapp/api/helpers"

	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
)

type File struct {
	File multipart.File `json:"file,omitempty" validate:"required"`
}

type Image struct {
	gorm.Model
	Url string `json:"url,omitempty" validate:"required"`
}
type ProfileImage struct {
	gorm.Model
	UserID  uint  `gorm:"primary_key column:user_id not null" json:"userId"`
	ImageID uint  `gorm:"primary_key column:image_id not null" json:"imageId"`
	Image   Image `json:"image"`
}

type PostImage struct {
	gorm.Model
	PostID  uint  `gorm:"primary_key column:post_id not null" json:"userId"`
	ImageID uint  `gorm:"primary_key column:image_id not null" json:"imageId"`
	Image   Image `json:"image"`
}

func (img *Image) Prepare() {
	img.ID = 0
	img.Url = html.EscapeString(strings.TrimSpace(img.Url))
}

func (img *ProfileImage) Prepare() {
	img.ID = 0
	img.ImageID = 0
	img.Image = Image{}
	img.UserID = 0
}

func (img *PostImage) Prepare() {
	img.ID = 0
	img.ImageID = 0
	img.Image = Image{}
	img.PostID = 0
}

var (
	validate = validator.New()
)

type mediaUpload interface {
	FileUpload(file File) (string, error)
}

type media struct{}

func NewMediaUpload() mediaUpload {
	return &media{}
}

func (*media) FileUpload(file File) (string, error) {
	err := validate.Struct(file)
	if err != nil {
		return "", err
	}
	uploadUrl, err := helpers.ImageUploadHelper(file.File)
	if err != nil {
		return "", err
	}
	if err != nil {
		return " ", err
	}
	return uploadUrl, nil
}

func (img *Image) SaveImage() (*Image, error) {
	err := db.Debug().Create(img).Error
	if err != nil {
		return &Image{}, err
	}
	return img, nil
}

func (img *ProfileImage) SaveProfileImage() (*ProfileImage, error) {
	err := db.Debug().Create(img).Error
	if err != nil {
		return &ProfileImage{}, err
	}
	return img, nil
}

func (img *PostImage) SavePostImage() (*PostImage, error) {
	err := db.Debug().Create(img).Error
	if err != nil {
		return &PostImage{}, err
	}
	return img, nil
}

func (img *Image) DeleteByID(imgid uint) (int64, error) {
	err := db.Debug().Table("images").Where("id = ? ", imgid).Take(&img).Delete(Image{})
	if err.Error != nil {
		return 0, err.Error
	}
	return db.RowsAffected, nil
}

func (img *Image) UpdateImageByID(imgid uint) (*Image, error) {
	db := GetDB().Debug().Table("images").Where("id=?", imgid).UpdateColumns(
		map[string]interface{}{
			"url": img.Url,
		},
	)
	if db.Error != nil {
		return &Image{}, db.Error
	}
	err := GetDB().Debug().Table("images").Where("id=?", imgid).Take(&img).Error
	if err != nil {
		return &Image{}, err
	}
	return img, nil
}

func (img *ProfileImage) DeleteProfileImgByUserID(uid uint) (int64, error) {
	err := db.Debug().Table("profile_images").Where("user_id = ? ", uid).Take(&img).Delete(Image{})
	if err.Error != nil {
		return 0, err.Error
	}
	return db.RowsAffected, nil
}

func (img *Image) DeletePostImgByID(imgid uint) (int64, error) {
	err := db.Debug().Table("post_images").Where("post_id = ? ", imgid).Take(&img).Delete(Image{})
	if err.Error != nil {
		return 0, err.Error
	}
	return db.RowsAffected, nil
}
