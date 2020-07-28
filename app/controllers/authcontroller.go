package controllers

import (
	"fmt"
	"github.com/gorilla/context"
	"gopattern/app/helpers"
	"gopattern/app/models"
	"net/http"
)

// GetAuthenticatedUser getting one user
func (app *App) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	user := &models.UserJSON{}

	userIDFromToken := fmt.Sprint(context.Get(r, "UserID"))
	userData, err := user.GetUser(userIDFromToken, app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	helpers.Success(w, http.StatusOK, "Hi " + userData.Name, userData)
	return
}
