package controllers

import (
	"encoding/json"
	"io/ioutil"

	"gopattern/app/helpers"
	"gopattern/app/models"
	"net/http"

	"github.com/gorilla/mux"
)

// GetAllRoles getting all users
func (app *App) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "Success", "Message": "Roles list"}
	role := &models.Role{}

	roles, err := role.GetRoles(app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	response["data"] = roles
	helpers.JSON(w, http.StatusOK, response)
}

// CreateRole create a new role
func (app *App) CreateRole(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "Success", "Message": "Role successfully created"}
	role := &models.Role{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &role)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Validate the role input
	err = role.Validate()
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	newRole, err := role.Create(app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	response["Data"] = newRole
	helpers.JSON(w, http.StatusCreated, response)
}

// GetRole get role by ID
func (app *App) GetRole(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "Success", "Message": "Role Detail"}
	role := &models.Role{}
	id := mux.Vars(r)["id"]

	roleData, err := role.GetRoleByID(id, app.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if roleData == nil {
		helpers.ERROR(w, http.StatusNotFound, err)
		return
	}

	response["Data"] = roleData
	helpers.JSON(w, http.StatusOK, response)
}

// Update Role
func (app *App) UpdateRole(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "Success", "Message": "Role updated"}
	role := &models.Role{}
	id := mux.Vars(r)["id"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &role)
	if err != nil {
		helpers.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Validate the role input
	err = role.Validate()
	if err != nil {
		helpers.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Check role data and update the data
	roleData, _ := role.GetRoleByID(id, app.DB)
	if roleData != nil {
		_, err := role.Update(id, app.DB)
		if err != nil {
			helpers.ERROR(w, http.StatusBadRequest, err)
			return
		}

		helpers.JSON(w, http.StatusOK, response)
		return
	}

	response["Status"] = "Error"
	response["Message"] = "Role not found"
	helpers.JSON(w, http.StatusNotFound, response)
	return
}

// DeleteRole delete selected role
func (app *App) DeleteRole(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"Status": "Success", "Message": "Role deleted"}
	role := &models.Role{}
	id := mux.Vars(r)["id"]

	roleData, _ := role.GetRoleByID(id, app.DB)
	if roleData != nil {
		_, err := role.Delete(roleData.ID, app.DB)
		if err != nil {
			helpers.ERROR(w, http.StatusBadRequest, err)
			return
		}
		helpers.JSON(w, http.StatusOK, response)
		return
	}

	response["Status"] = "Error"
	response["Message"] = "Role not found"
	helpers.JSON(w, http.StatusNotFound, response)
	return
}
