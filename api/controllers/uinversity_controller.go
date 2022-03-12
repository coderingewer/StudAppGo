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

func CreateUniversity(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	university := models.University{}
	err = json.Unmarshal(body, &university)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	university.Prepare()
	universityCreated, err := university.SAve()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	university.CityID = universityCreated.CityID
	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.URL, university.ID))
	utils.JSON(w, http.StatusCreated, universityCreated)
}

func GetByCityID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["cityId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	university := models.University{}
	universities, err := university.FindByCityID(uint(cid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, universities)
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	university := models.University{}
	universities, err := university.FindBYID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, universities)
}

func AddAFaculty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	fid, err := strconv.ParseUint(vars["facultyId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	universityFaculty := models.UniverstyFaculty{}
	unif, err := universityFaculty.AddAFacultyByID(uint(unid), uint(fid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, unif)
}
