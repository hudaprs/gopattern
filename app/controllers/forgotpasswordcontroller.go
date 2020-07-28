package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopattern/app/helpers"
	"gopattern/app/models"
	"io/ioutil"
	"net/http"
)

func (app App) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	verification := &models.Verification{}

	// Get all request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Unmarshal the request
	err = json.Unmarshal(body, &user)

	// Validate the user input
	err = user.Validate("forgot-password")
	if err != nil {
		helpers.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Get the user data by E-Mail
	userData, err := user.GetUserByEmail(app.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if userData != nil {
		// Check the verification data first
		userIDString := fmt.Sprint(userData.ID)
		verificationData, _ := verification.GetVerificationByID(userIDString, "Forgot Password", app.DB)
		if verificationData != nil {
			// Delete the existing verification data
			verificationIDString := fmt.Sprint(verificationData.ID)
			_, err := verification.DeleteVerification(verificationIDString, app.DB)
			if err != nil {
				helpers.Error(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		// Generate the new verification data
		randomString := helpers.RandStringRunes(30)
		verification.Name = "Forgot Password"
		verification.Token = randomString
		verification.UserID = userData.ID
		app.DB.Save(&verification)

		helpers.Success(w, http.StatusCreated, "Verification token has been sent", verification)
		return
	}

	helpers.Error(w, http.StatusNotFound, "User not found")
	return
}

// ChangePassword user
func (app *App) ChangePassword(w http.ResponseWriter, r *http.Request) {
	verification := &models.Verification{}
	user := &models.User{}
	userJSON := &models.UserJSON{}
	verificationToken := mux.Vars(r)["token"]

	// Check verification token
	if verificationToken == "" {
		helpers.Error(w, http.StatusNotFound, "Verification token not found")
		return
	}

	// Get verification data by token
	verificationData, _ := verification.GetVerificationByToken(verificationToken, app.DB)
	if verificationData == nil {
		helpers.Error(w, http.StatusNotFound, "Verification data not found")
		return
	}

	// Get all request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Unmarshal the request
	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the user input
	err = user.Validate("change-password")
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get the user data
	userIDString := fmt.Sprint(verificationData.UserID)
	userData, _ := userJSON.GetUser(userIDString, app.DB)
	if userData == nil {
		helpers.Error(w, http.StatusNotFound, "User not found")
		return
	}

	// Update the user password
	_, err = user.ChangeUserPassword(userIDString, app.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Delete verification data
	verificationIDString := fmt.Sprint(verificationData.ID)
	_, err = verification.DeleteVerification(verificationIDString, app.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.Success(w, http.StatusOK, "Password successfully changed", userData)
	return
}
