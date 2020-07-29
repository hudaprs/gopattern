package controllers

import (
	"fmt"
	"github.com/gorilla/context"
	"gopattern/app/helpers"
	"gopattern/app/models"
	"gopattern/config"
	"net/http"
)


// GetAuthenticatedUser getting one user
func GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	user := &models.UserJSON{}

	userIDFromToken := fmt.Sprint(context.Get(r, "UserID"))
	userData, err := user.GetUser(userIDFromToken, config.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.Success(w, http.StatusOK, "Hi " + userData.Name, userData)
	return
}
