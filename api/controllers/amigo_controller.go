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

func DeleteAmigo(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	amigo := models.Amigo{}
	err = models.GetDB().Debug().Table("amigos").Where("user_id = ?", uid).Take(&amigo).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if uint(uid) != amigo.UserID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	_, err = amigo.DeleteAmigoByUserID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	utils.JSON(w, http.StatusNoContent, " ")
}

func GetAmigosByDESC(w http.ResponseWriter, r *http.Request) {
	amigo := models.Amigo{}

	amigos, err := amigo.FindAllAmigosByDESC()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, amigos)
}

func GetAmigosByUserID(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	amigo := models.Amigo{}
	amigos, err := amigo.FindByUserID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, amigos)
}

func GetAmigosByCityID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["cityId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	amigo := models.Amigo{}
	amigos, err := amigo.FindByCityID(uint(cid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, amigos)
}
