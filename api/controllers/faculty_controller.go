package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"studapp/api/models"
	"studapp/api/utils"

	"github.com/gorilla/mux"
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

func GetFacultyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	faculty := models.Faculty{}
	facultyReceived, err := faculty.FindFacultyByID(uint(fid))
	utils.JSON(w, http.StatusOK, facultyReceived)
}

func GetFacultyByUniID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	faculty := models.UniversityFaculty{}
	faculties, err := faculty.GetFacultyByUniID(uint(unid))
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	utils.JSON(w, http.StatusOK, faculties)
}

func DeleteFacultyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	faculty := models.Faculty{}
	_, err = faculty.DeleteByID(uint(fid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", fid))
	utils.JSON(w, http.StatusNoContent, "")
}
