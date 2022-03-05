package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"studapp/api/auth"
	"studapp/api/models"
	"studapp/api/utils"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		fmt.Println(string(err.Error()))
		return
	}
	utils.JSON(w, http.StatusOK, token)
}

func SignIn(email, password string) (string, error) {
	var err error
	user := models.User{}
	err = models.GetDB().Debug().Table("users").Where("email=?", email).First(&user).Error
	if err != nil {
		fmt.Println(string(err.Error()))
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
