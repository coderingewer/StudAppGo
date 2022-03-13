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

func AddADepartment(w http.ResponseWriter, r *http.Request) {
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
	did, err := strconv.ParseUint(vars["departmentId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	universityDepartment := models.UniversityDepartment{}
	unif, err := universityDepartment.AddADepartmentByID(uint(unid), uint(did), uint(fid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, unif)
}

func GetByDepartmentID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, err := strconv.ParseUint(vars["departmentId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	university := models.UniversityDepartment{}
	universities, err := university.FindUniByDepartmentID(uint(did))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, universities)
}

func GetByFacultyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, err := strconv.ParseUint(vars["facultyId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	university := models.UniverstyFaculty{}
	universities, err := university.FindUniByFacultyID(uint(did))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, universities)
}

func DeleteUniversity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unid, err := strconv.ParseUint(vars["id"], 10, 65)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	university := models.University{}
	_, err = university.DeleteByID(uint(unid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusNoContent, "")
}

func DeleteUniDepartmentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unidid, err := strconv.ParseUint(vars["id"], 10, 65)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	department := models.UniversityDepartment{}
	_, err = department.DeleteUniversityDepartmentByUniID(uint(unidid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusNoContent, "")
}

func DeleteUniFacultyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unifid, err := strconv.ParseUint(vars["id"], 10, 65)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	faculty := models.UniverstyFaculty{}
	_, err = faculty.DeleteUniversityFacultyByID(uint(unifid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusNoContent, "")
}
