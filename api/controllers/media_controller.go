package controllers

import (
	"net/http"
	"strconv"
	"studapp/api/auth"
	"studapp/api/models"
	"studapp/api/utils"

	"github.com/gorilla/mux"
)

func ImgUpload(w http.ResponseWriter, r *http.Request) {
	img := models.Image{}
	formFile, _, err := r.FormFile("file")
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	img.Url = uploadUrl
	image, err := img.SaveImage()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, image)
}

func UploadProfileImg(w http.ResponseWriter, r *http.Request) {
	img := models.Image{}
	profileImage := models.ProfileImage{}
	formFile, _, err := r.FormFile("file")
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	profileImage.UserID = uid
	profileImage.ImageID = img.ID

	prflImg, err := profileImage.SaveProfileImage()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	img.Url = uploadUrl
	_, err = img.SaveImage()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, prflImg)
}

func UpdateImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imgid, err := strconv.ParseUint(vars["imageId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	formFile, _, err := r.FormFile("file")
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	img := models.Image{}

	err = models.GetDB().Debug().Table("images").Where("id = ?", imgid).Take(&img).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	imageUpdate := models.Image{}

	imageUpdate.Prepare()
	imageUpdate.ID = img.ID
	imageUpdate.Url = uploadUrl

	imgUpdated, err := imageUpdate.UpdateImageByID(uint(imgid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, imgUpdated)

}

func UpdateProfileImage(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	formFile, _, err := r.FormFile("file")
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	img := models.ProfileImage{}

	err = models.GetDB().Debug().Table("profile_images").Where("user_id = ?", uid).Take(&img).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	image := models.Image{}

	err = models.GetDB().Debug().Table("images").Where("id = ?", img.ImageID).Take(&image).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	imageUpdate := models.Image{}
	imageUpdate.Prepare()
	imageUpdate.ID = img.ID
	imageUpdate.Url = uploadUrl
	profileImgUpdated, err := image.UpdateImageByID(uint(img.ImageID))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, profileImgUpdated)
}

func DeleteImageByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imgid, err := strconv.ParseUint(vars["imageId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	img := models.Image{}
	_, err = img.DeleteByID(uint(imgid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusNoContent, "")

}
