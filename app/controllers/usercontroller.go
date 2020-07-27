package controllers

import (
	"encoding/json"
	"fmt"
	"gopattern/app/helpers"
	"gopattern/app/models"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/context"
)

// Register a new user
func (app *App) Register(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "Success", "Message": "Register succesfully"}
	user := &models.User{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Validate the user
	err = user.Validate("register")
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Check the user
	checkUser, _ := user.GetUserByEmail(app.DB)
	if checkUser != nil {
		response["Status"] = "Error"
		response["Message"] = "User already registered"
		helpers.JSON(w, http.StatusUnauthorized, response)
		return
	}

	_, err = user.Register(app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	helpers.JSON(w, http.StatusCreated, response)
	return
}

// Login a user
func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "Success", "Message": "Login Success"}
	user := &models.User{}
	role := &models.Role{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Validate user
	err = user.Validate("login")
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Check the user
	checkUser, _ := user.GetUserByEmail(app.DB)
	if checkUser != nil {
		// Check Password Hash
		err = user.CheckHashedPassword(checkUser.Password, user.Password)
		if err != nil {
			helpers.ERROR(w, http.StatusBadRequest, err)
			return
		}

		// Get Role for user
		makeIDtoString := fmt.Sprint(checkUser.RoleID)
		role, err := role.GetRoleByID(makeIDtoString, app.DB)
		if err != nil {
			helpers.ERROR(w, http.StatusBadRequest, err)
			return
		}

		token, err := helpers.EncodeAuthToken(checkUser.ID, checkUser.Name, checkUser.Email, role.Name)
		if err != nil {
			helpers.ERROR(w, http.StatusBadRequest, err)
			return
		}
		response["Token"] = token
		helpers.JSON(w, http.StatusOK, response)
		return
	}

	response["Status"] = "Error"
	response["Message"] = "User not found"
	helpers.JSON(w, http.StatusNotFound, response)
	return
}

// GetAllUsers getting all users
func (app *App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "Success", "Message": "Users List"}
	user := &models.UserJSON{}

	users, err := user.GetUsers(app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	response["Data"] = users
	helpers.JSON(w, http.StatusOK, response)

	return
}

// GetOneUser getting one user
func (app *App) GetOneUser(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "Success", "Message": "User Detail"}
	user := &models.UserJSON{}

	userIDFromToken := fmt.Sprint(context.Get(r,"UserID"))
	userData, err := user.GetUser(userIDFromToken, app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	response["Data"] = userData
	helpers.JSON(w, http.StatusOK, response)
	return
}