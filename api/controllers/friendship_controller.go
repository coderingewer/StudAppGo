package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"studapp/api/auth"
	"studapp/api/models"
	"studapp/api/utils"

	"github.com/gorilla/mux"
)

func SendRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rid, err := strconv.ParseUint(vars["userId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	fr := models.FriendshipRequest{}
	fr.Prepare()
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	fr.SenderID = uint(uid)
	fr.RecieverID = uint(rid)
	friendshipRequested, err := fr.CreateRequest()
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	utils.JSON(w, http.StatusOK, friendshipRequested)
}

func AcceptFrienshipRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	frid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	fr := models.FriendshipRequest{}
	err = models.GetDB().Debug().Table("friendship_requests").Where("reciever_id = ?", uid).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if uid != fr.RecieverID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	userf := models.UserFriend{}
	userf.Prepare()
	userf.SenderID = fr.SenderID
	userf.RecieverID = uint(uid)

	FriendCreated, err := userf.CreateFriend()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
	_, err = fr.DeleteFrienshipRequestByUserID(uint(frid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}

	utils.JSON(w, http.StatusOK, FriendCreated)
}

func GetRequestsByRecieverID(w http.ResponseWriter, r *http.Request) {
	reciever := models.FriendshipRequest{}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	friendRequests, err := reciever.FindRequestsByRecieverID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	utils.JSON(w, http.StatusOK, friendRequests)

}

func GetRequestsBySenderID(w http.ResponseWriter, r *http.Request) {
	reciever := models.FriendshipRequest{}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	friendRequests, err := reciever.FindRequestsBySenderID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	utils.JSON(w, http.StatusOK, friendRequests)

}

func GetFriendsByUserID(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	userf := models.UserFriend{}
	friends, err := userf.FindFriendsByUserID(uid)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}

	utils.JSON(w, http.StatusOK, friends)
}

func DeleteRequestBySenderID(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	fr := models.FriendshipRequest{}
	err = models.GetDB().Debug().Table("friendship_requests").Where("sender_id = ?", uid).Take(fr).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if fr.SenderID != uid {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	_, err = fr.DeleteFrienshipRequestByUserID(uid)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	utils.JSON(w, http.StatusNoContent, "")
}

func DeleteRequestByRecieverID(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	fr := models.FriendshipRequest{}
	err = models.GetDB().Debug().Table("friendship_requests").Where("reciever_id").Take(fr).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if fr.RecieverID != uid {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	_, err = fr.DeleteFrienshipRequestByUserID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	utils.JSON(w, http.StatusNoContent, "")
}
