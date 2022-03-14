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

func CreateDepartment(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	department := models.Department{}
	err = json.Unmarshal(body, &department)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	department.Prepare()
	departmentCreated, err := department.Save()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusCreated, departmentCreated)
}

func GetDepartments(w http.ResponseWriter, r *http.Request) {
	department := models.Department{}
	faculties, err := department.FindAllDepartments()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, faculties)
}

func GetDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	department := models.Department{}
	departmentReceived, err := department.FindDepartmentByID(uint(did))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, departmentReceived)
}

func DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	department := models.Department{}
	err = models.GetDB().Debug().Table("departments").Where("id = ?", did).Take(&department).Error
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		return
	}

	_, err = department.DeleteByID(uint(did))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", did))
	utils.JSON(w, http.StatusNoContent, "")
}

func GetDepartmentByUniID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unid, err := strconv.ParseUint(vars["universityId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	department := models.UniversityDepartment{}
	departments, err := department.FindDepartmentByUniID(uint(unid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, departments)
}

func GetDepartmentByUniIDAndFacultyID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unid, err := strconv.ParseUint(vars["universityId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	fid, err := strconv.ParseUint(vars["facultyId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	department := models.UniversityDepartment{}
	departments, err := department.FindDepartmentByFacultyIDAndUniID(uint(unid), uint(fid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, departments)
}

func GetDepartmentByFacultyID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["FacultyId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	department := models.UniversityDepartment{}
	universities, err := department.FindDepartmentByFacultyID(uint(fid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, universities)
}

func DeleteDepartmentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	department := models.Department{}
	_, err = department.DeleteByID(uint(fid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", fid))
	utils.JSON(w, http.StatusNoContent, "")
}
