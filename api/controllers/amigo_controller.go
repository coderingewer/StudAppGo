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

func CreateAmigo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	amigo := models.Amigo{}
	err = json.Unmarshal(body, &amigo)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	amigo.Prepare()

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	}

	amigo.UserID = uint(uid)
	amigoCreated, err := amigo.CreateAmigo()
	if err != nil {
		formatErrror := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formatErrror)
		w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.URL, amigo.ID))
		utils.JSON(w, http.StatusCreated, amigoCreated)
	}
}

func DelteAmigo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	aid, err := strconv.ParseUint(vars["amigoId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	amigo := models.Amigo{}
	err = models.GetDB().Debug().Table("amigos").Where("id = ?", aid).Take(&amigo).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	if uint(uid) != amigo.UserID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return

	}
	_, err = amigo.DeleteAmigoByID(uint(aid))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", aid))
	utils.JSON(w, http.StatusNoContent, " ")
}
