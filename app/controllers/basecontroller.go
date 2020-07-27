package controllers

import (
	"fmt"
	"gopattern/app/middlewares"
	"gopattern/app/models"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres
	"github.com/joho/godotenv"
)

// App Struct
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Routes app
func (app *App) Routes() {

	app.Router = mux.NewRouter().StrictSlash(true)
	ProtectedRouter := app.Router.PathPrefix("/api/v1").Subrouter()
	ProtectedRouterHighAdmin := app.Router.PathPrefix("/api/v1").Subrouter()

	// Middlewares
	app.Router.Use(middlewares.SetContentTypeHeader)
	ProtectedRouter.Use(middlewares.AuthJwtVerify)
	ProtectedRouterHighAdmin.Use(middlewares.AuthJwtVerify)
	ProtectedRouterHighAdmin.Use(middlewares.OnlyHighAdmin)

	// Open Routes
	app.Router.HandleFunc("/api/register", app.Register).Methods("POST")
	app.Router.HandleFunc("/api/login", app.Login).Methods("POST")

	// High Admin Routes
	ProtectedRouterHighAdmin.HandleFunc("/roles", app.GetAllRoles).Methods("GET")
	ProtectedRouterHighAdmin.HandleFunc("/roles", app.CreateRole).Methods("POST")
	ProtectedRouterHighAdmin.HandleFunc("/roles/{id}", app.GetRole).Methods("GET")
	ProtectedRouterHighAdmin.HandleFunc("/roles/{id}", app.UpdateRole).Methods("PATCH")
	ProtectedRouterHighAdmin.HandleFunc("/roles/{id}", app.DeleteRole).Methods("DELETE")
	ProtectedRouterHighAdmin.HandleFunc("/users", app.GetAllUsers).Methods("GET")

	// Protected Router
	ProtectedRouter.HandleFunc("/users/me", app.GetOneUser).Methods("GET")
	ProtectedRouter.HandleFunc("/users/me/upload-image", app.UploadUserImage).Methods("PUT")
	ProtectedRouter.HandleFunc("/users/me/delete-image", app.DeleteImage).Methods("DELETE")
}

// Init App
func (app *App) Init(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	var err error
	DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s	", DbHost, DbPort, DbUser, DbName, DbPassword)

	app.DB, err = gorm.Open("postgres", DBURI)
	if err != nil {
		fmt.Println("Failed connecting to database")
		panic(err)
	}

	fmt.Println("Connected To Database")
	fmt.Println("Server started port 8000")

	app.DB.Debug().AutoMigrate(&models.User{}, &models.Role{})
	app.Routes()

	log.Fatal(http.ListenAndServe(":8000", app.Router))
}

// RunServer Run App Server
func (app *App) RunServer() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load env")
		panic(err)
	}

	app.Init(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)
}
