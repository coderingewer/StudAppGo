package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"studapp/api/models"
	"studapp/api/utils"
)

func CreateFaculty(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	faculty := models.Faculty{}
	err = json.Unmarshal(body, &faculty)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	faculty.Prepare()
	facultyCreated, err := faculty.Save()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusCreated, facultyCreated)
}

func GetFaculties(w http.ResponseWriter, r *http.Request) {
	faculty := models.Faculty{}
	faculties, err := faculty.FindAllFaculty()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, faculties)
}
