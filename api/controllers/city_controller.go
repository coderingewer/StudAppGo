package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"studapp/api/models"
	"studapp/api/utils"
)

func CreateCity(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	city := models.City{}

	err = json.Unmarshal(body, &city)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	city.Prepare()
	cityCreated, err := city.Save()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.URL, city.ID))
	utils.JSON(w, http.StatusCreated, cityCreated)
}

func GetAllCities(w http.ResponseWriter, r *http.Request) {
	city := models.City{}

	cities, err := city.FindAll()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	utils.JSON(w, http.StatusOK, cities)
}
