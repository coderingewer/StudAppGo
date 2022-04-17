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
	img.Url = uploadUrl
	_, err = img.SaveImage()
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
	utils.JSON(w, http.StatusOK, prflImg)
}

func UploadPostImg(w http.ResponseWriter, r *http.Request) {
	img := models.Image{}
	postImage := models.PostImage{}
	formFile, _, err := r.FormFile("file")
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["postId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
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
	postImage.PostID = uint(pid)
	postImage.ImageID = img.ID
	prflImg, err := postImage.SavePostImage()
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

func UpdatePostImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["postId"], 10, 64)
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
	img := models.PostImage{}

	err = models.GetDB().Debug().Table("post_images").Where("post_id = ?", pid).Take(&img).Error
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
	imageUpdate.ID = image.ID
	imageUpdate.Url = uploadUrl
	profileImgUpdated, err := imageUpdate.UpdateImageByID(uint(img.ImageID))
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

func DeleteImageByUserID(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	img := models.ProfileImage{}
	_, err = img.DeleteProfileImgByUserID(uid)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusNoContent, "")
}
