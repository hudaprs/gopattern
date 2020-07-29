package controllers

import (
	"fmt"
	"gopattern/app/helpers"
	"gopattern/app/models"
	"gopattern/config"
	"net/http"

	"github.com/gorilla/context"
)

// GetAuthenticatedUser getting one user
func GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	user := &models.UserJSON{}

	userIDFromToken := fmt.Sprint(context.Get(r, "UserID"))
	userData, _ := user.GetUser(userIDFromToken, config.DB)
	if userData == nil {
		helpers.Error(w, http.StatusBadRequest, "User not found")
		return
	}

	helpers.Success(w, http.StatusOK, "Hi "+userData.Name, userData)
	return
}
