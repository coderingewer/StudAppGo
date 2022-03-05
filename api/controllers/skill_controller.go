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

func CreateSkill(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	skill := models.Skill{}
	err = json.Unmarshal(body, &skill)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	skill.Prepare()
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	skill.AuthorID = uid
	skillCreated, err := skill.Save()
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.URL, skill.ID))
	utils.JSON(w, http.StatusCreated, skillCreated)
}

func GetSkills(w http.ResponseWriter, r *http.Request) {
	skill := models.Skill{}
	skills, err := skill.FindAllSkills()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, skills)
}

func GetSkill(w http.ResponseWriter, r *http.Request) {
	skill := models.Skill{}
	vars := mux.Vars(r)
	sid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	skillRecieved, err := skill.FindByID(uint(sid))
	utils.JSON(w, http.StatusOK, skillRecieved)
}

func UpdateSkill(w http.ResponseWriter, r *http.Request) {
	skill := models.Skill{}
	vars := mux.Vars(r)
	sid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
	}
	err = models.GetDB().Debug().Table("skills").Where("id = ?", sid).Take(&skill).Error
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		return
	}
	if uid != skill.AuthorID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	skillUpdate := models.Skill{}
	err = json.Unmarshal(body, &skillUpdate)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	skillUpdate.AuthorID = skill.AuthorID
	if skillUpdate.AuthorID != uid {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
	}
	skillUpdate.Prepare()
	skillUpdate.ID = skill.ID
	skillUpdated, err := skillUpdate.UpdateSkill(uint(sid))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	utils.JSON(w, http.StatusOK, skillUpdated)
}

func DeleteSkill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}

	skill := models.Skill{}
	err = models.GetDB().Debug().Table("skills").Where("id = ?", sid).Take(&skill).Error
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		return
	}
	if skill.AuthorID != uint(uid) {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
	}
	_, err = skill.DeleteByID(uint(sid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusNoContent, "")
}

func GetSkillsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	skill := models.Skill{}
	uid, err := strconv.ParseUint(vars["userId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	skills, err := skill.FindByUserID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, skills)
}
