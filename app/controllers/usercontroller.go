package controllers

import (
	"encoding/json"
	"fmt"
	"gopattern/app/helpers"
	"gopattern/app/models"
	"gopattern/config"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/context"
)

// Register a new user
func Register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the user
	err = user.ValidateRegister(config.DB)
	if err != nil {
		helpers.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Check the user
	checkUser, _ := user.GetUserByEmail(config.DB)
	if checkUser != nil {
		helpers.Error(w, http.StatusUnauthorized, "Email already registered")
		return
	}

	userData, err := user.Register(config.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	data := map[string]interface{}{
		"ID":        userData.ID,
		"Name":      userData.Name,
		"Email":     userData.Email,
		"RoleID":    userData.RoleID,
		"CreatedAt": userData.CreatedAt,
		"UpdatedAt": userData.CreatedAt,
		"DeletedAt": userData.DeletedAt,
	}
	helpers.Success(w, http.StatusCreated, "User successfully registered", data)
	return
}

// Login a user
func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	role := &models.Role{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate user
	err = user.Validate("login")
	if err != nil {
		helpers.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Check the user
	userData, _ := user.GetUserByEmail(config.DB)
	if userData != nil {
		// Check Password Hash
		err = user.CheckHashedPassword(userData.Password, user.Password)
		if err != nil {
			helpers.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Get Role for user
		makeIDtoString := fmt.Sprint(userData.RoleID)
		role, _ := role.GetRoleByID(makeIDtoString, config.DB)
		if role == nil {
			helpers.Error(w, http.StatusBadRequest, "Role data not found")
			return
		}

		token, err := helpers.EncodeAuthToken(userData.ID, userData.Name, userData.Email, role.Name)
		if err != nil {
			helpers.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		mapUserData := map[string]interface{}{
			"ID":        userData.ID,
			"Name":      userData.Name,
			"Email":     userData.Email,
			"Role":      userData.Role,
			"RoleID":    userData.RoleID,
			"CreatedAt": userData.CreatedAt,
			"UpdatedAt": userData.CreatedAt,
			"DeletedAt": userData.DeletedAt,
		}

		data := map[string]interface{}{"Token": token, "User": mapUserData}
		helpers.Success(w, http.StatusOK, "Login successfully", data)
		return
	}

	helpers.Error(w, http.StatusNotFound, "Invalid credentials")
	return
}

// GetAllUsers getting all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	user := &models.UserJSON{}

	// Get total of user data
	total, err := user.CountUsers(config.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Paginate the users
	queryParams := r.URL.Query()
	nameParam := queryParams.Get("name")
	limitParam, _ := strconv.Atoi(queryParams.Get("limit"))
	if limitParam < 1 {
		limitParam = 10
	}
	pages := total / limitParam
	if (total % limitParam) != 0 {
		pages++
	}

	// Get the pagination data
	page, begin := helpers.Pagination(r, limitParam)

	users, err := user.GetUsers(begin, page, nameParam, config.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	mapUsers := helpers.PaginationResponse(r, page, pages, limitParam, total, users)

	helpers.Success(w, http.StatusOK, "Users list", mapUsers)
	return
}

// UploadUserImage to local
func UploadUserImage(w http.ResponseWriter, r *http.Request) {
	// Update the user
	user := &models.UserJSON{}
	userIDFromContext := fmt.Sprint(context.Get(r, "UserID"))
	if err := config.DB.Debug().Table("users").Preload("Role").Where("id = ?", userIDFromContext).First(&user).Error; err != nil {
		helpers.Error(w, http.StatusBadRequest, "User not found")
		return
	}

	// Parse input to type multipart/form-data
	// Set the maximum file size
	r.ParseMultipartForm(10 << 20)

	// Retreive file from posted form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	if file == nil {
		helpers.Error(w, http.StatusBadRequest, "File is required")
		return
	}

	// Get header type
	headerType := handler.Header["Content-Type"][0]
	headerTypesArray := []string{"image/png", "image/jpeg", "image/jpg"}
	headerTypes := map[string]string{}
	for _, header := range headerTypesArray {
		headerTypes[header] = header
	}

	// Check the type of the file
	if headerType != headerTypes[headerType] {
		helpers.Error(w, http.StatusUnprocessableEntity, "The file must be png, jpeg, or jpg")
		return
	}

	// Write temporary file in local
	getFileExtension := strings.Split(headerType, "/")[1]
	tempFile, err := ioutil.TempFile("static/user_images", "images-*."+getFileExtension)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	defer tempFile.Close()

	// Get The file bytes
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	// Write the file
	tempFile.Write(fileBytes)

	// Remove Previous image (if exists)
	if user.ImageURL != "" {
		getFileNameOnly := strings.Split(user.ImageURL, "/")[3]
		err := os.Remove("static/user_images/" + getFileNameOnly)
		if err != nil {
			helpers.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		user.ImageURL = r.Host + "/" + strings.ReplaceAll(tempFile.Name(), "\\", "/")
	} else {
		// Upload file as usual
		user.ImageURL = r.Host + "/" + strings.ReplaceAll(tempFile.Name(), "\\", "/")
	}

	// Save the user
	config.DB.Save(&user)

	helpers.Success(w, http.StatusOK, "Image uploaded", user)
}

// DeleteImage user
func DeleteImage(w http.ResponseWriter, r *http.Request) {
	user := &models.UserJSON{}
	userIDFromContext := fmt.Sprint(context.Get(r, "UserID"))

	// Get One User
	userData, err := user.GetUser(userIDFromContext, config.DB)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if the user isn't nil, and remove the image
	if userData != nil {
		// Check if user didn't have any image
		if userData.ImageURL == "" {
			helpers.Error(w, http.StatusBadRequest, "User din't have image, yet")
			return
		}

		getFileNameOnly := strings.Split(userData.ImageURL, "/")[3]
		err := os.Remove("static/user_images/" + getFileNameOnly)
		if err != nil {
			helpers.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		// Set
		userData.ImageURL = ""
		config.DB.Save(&userData)

		helpers.Success(w, http.StatusOK, "Image deleted", userData)
		return
	}

	helpers.Error(w, http.StatusNotFound, "User not found")
	return
}
