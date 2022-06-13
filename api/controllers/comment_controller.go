package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"studapp/api/auth"
	"studapp/api/models"
	"studapp/api/utils"

	"github.com/gorilla/mux"
)

func CreateCommet(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["postId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	cmmt := models.Comment{}
	err = json.Unmarshal(body, &cmmt)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	cmmt.Prepare()
	cmmt.UserID = uid
	cmmt.PostID = uint(pid)
	commentCreated, err := cmmt.CreateCommet()
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.URL, cmmt.ID))
	utils.JSON(w, http.StatusCreated, commentCreated)
}

func GetByUserID(w http.ResponseWriter, r *http.Request) {
	comment := models.Comment{}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	Comments, err := comment.FindByUserID(uid)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, Comments)
}

func GetByPostID(w http.ResponseWriter, r *http.Request) {

	comment := models.Comment{}
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["postId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	Comments, err := comment.FindByUserID(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, Comments)
}

func DeleteCommentByPostID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["postId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	comment := models.Comment{}

	err = models.GetDB().Debug().Table("comments").Where("post_id", pid).Take(&comment).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if comment.PostID != uint(pid) {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if comment.UserID != uid {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}

	_, err = comment.DeleteCommetByPostID(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, "")
}
