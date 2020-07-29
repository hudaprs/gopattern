package controllers

import (
	"encoding/json"
	"gopattern/config"
	"io/ioutil"
	"strconv"

	"gopattern/app/helpers"
	"gopattern/app/models"
	"net/http"

	"github.com/gorilla/mux"
)

// GetAllRoles getting all users
func GetAllRoles(w http.ResponseWriter, r *http.Request) {
	role := &models.Role{}

	// Count total of roles
	total, err := role.CountRoles(config.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Paginate the roles
	queryParams := r.URL.Query()
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	nameParam := queryParams.Get("name")
	if limit < 1 {
		limit = 10
	}
	page, begin := helpers.Pagination(r, limit)
	pages := total / limit
	if (total % limit) != 0 {
		pages++
	}

	// Return the paginate
	roles, err := role.GetRoles(begin, limit, nameParam, config.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	mapRoles := helpers.PaginationResponse(r, page, pages, limit, total, roles)

	helpers.Success(w, http.StatusOK, "Roles list", mapRoles)
	return
}

// CreateRole create a new role
func CreateRole(w http.ResponseWriter, r *http.Request) {
	role := &models.Role{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(body, &role)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the role input
	err = role.Validate()
	if err != nil {
		helpers.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	newRole, err := role.Create(config.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.Success(w, http.StatusCreated, "Role successfully created", newRole)
}

// GetRole get role by ID
func GetRole(w http.ResponseWriter, r *http.Request) {
	role := &models.Role{}
	id := mux.Vars(r)["id"]

	roleData, _ := role.GetRoleByID(id, config.DB)

	if roleData == nil {
		helpers.Error(w, http.StatusNotFound, "Role not found")
		return
	}

	helpers.Success(w, http.StatusOK, "Role Detail", roleData)
}

// Update Role
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	role := &models.Role{}
	id := mux.Vars(r)["id"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(body, &role)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the role input
	err = role.Validate()
	if err != nil {
		helpers.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Check role data and update the data
	roleData, _ := role.GetRoleByID(id, config.DB)
	if roleData != nil {
		if err := config.DB.Debug().Table("roles").First(&roleData).Update("name", role.Name).Error; err != nil {
			helpers.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		config.DB.Save(&roleData)

		helpers.Success(w, http.StatusOK, "Role successfully updated", roleData)
		return
	}

	helpers.Error(w, http.StatusNotFound, "Role not found")
	return
}

// DeleteRole delete selected role
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	role := &models.Role{}
	id := mux.Vars(r)["id"]

	roleData, _ := role.GetRoleByID(id, config.DB)
	if roleData != nil {
		_, err := role.Delete(roleData.ID, config.DB)
		if err != nil {
			helpers.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		helpers.Success(w, http.StatusOK, "Role successfully deleted", roleData)
		return
	}

	helpers.Error(w, http.StatusNotFound, "Role not found")
	return
}
