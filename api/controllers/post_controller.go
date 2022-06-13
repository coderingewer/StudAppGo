package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"studapp/api/auth"
	"studapp/api/models"
	"studapp/api/utils"

	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post.Prepare()
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	post.UserID = uint(uid)
	postCreated, err := post.Save()
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.URL, post.ID))
	utils.JSON(w, http.StatusCreated, postCreated)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	posts, err := post.FindAllPosts()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
	}
	post := models.Post{}
	postReceived, err := post.FindByID(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, postReceived)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	post := models.Post{}
	err = models.GetDB().Debug().Table("posts").Where("id = ?", pid).Take(&post).Error
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		return
	}

	if uid != post.UserID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	postUpdate.UserID = post.UserID
	if uid != postUpdate.UserID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	postUpdate.Prepare()
	postUpdate.ID = post.ID

	postUpdated, err := postUpdate.UpdatePost(uint(pid))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	utils.JSON(w, http.StatusOK, postUpdated)

}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	post := models.Post{}
	err = models.GetDB().Debug().Table("posts").Where("id = ?", pid).Take(&post).Error
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		return
	}

	if uid != post.UserID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	img := models.PostImage{}

	_, err = img.Image.DeletePostImgByID(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	_, err = post.DeleteByID(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = img.Image.DeletePostImgByID(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	utils.JSON(w, http.StatusNoContent, "")
}

func GetPostsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post := models.Post{}
	uid, err := strconv.ParseUint(vars["userId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	posts, err := post.FinBYUserID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, posts)
}
