package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopattern/app/models"
)

var DB *gorm.DB

func Connect(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	var err error
	DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s	", DbHost, DbPort, DbUser, DbName, DbPassword)

	DB, err = gorm.Open("postgres", DBURI)
	if err != nil {
		fmt.Println("Failed connecting to database")
		panic(err)
	}
	// Migrate the models
	DB.Debug().AutoMigrate(&models.User{}, &models.Role{}, &models.Verification{})
}
