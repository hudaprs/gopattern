package models

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

// User Struct
type User struct {
	gorm.Model
	Name     string `gorm:"size:100;not null;"`
	Email    string `gorm:"size:100;not null;"`
	Password string `gorm:"size:255;not null;"`
	Role     Role   `gorm:"ForeignKey:RoleID"`
	RoleID   uint   `gorm:"not null"`
	ImageURL string `gorm:"size:255"`
}

// UserJSON struct
type UserJSON struct {
	gorm.Model
	Name     string `gorm:"size:100;not null;"`
	Email    string `gorm:"size:100;not null;"`
	Role     Role   `gorm:"ForeignKey:RoleID"`
	RoleID   uint   `gorm:"not null"`
	ImageURL string `gorm:"size:255"`
}

// Set User's table name to be `profiles`
func (UserJSON) TableName() string {
	return "users"
}

// HashPassword of user
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckHashedPassword of user
func (user *User) CheckHashedPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password))
	if err != nil {
		return errors.New("Invalid Credentials")
	}
	return nil
}

// BeforeSave user password must be hashed
func (user *User) BeforeSave() error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return nil
}

// Validate user
func (user User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "register":
		if user.Name == "" {
			return errors.New("Name is required")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Email is required and must be valid")
		}
		if user.Password == "" {
			return errors.New("Password is required")
		}
		if user.RoleID == 0 {
			return errors.New("Role ID is required")
		}
		return nil
	case "login":
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Email must be valid")
		}
		if user.Password == "" {
			return errors.New("Password is required")
		}
		return nil
	case "forgot-password":
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Email is required and must be valid")
		}
		return nil
	case "change-password":
		if user.Password == "" {
			return errors.New("New password is required")
		}
		return nil
	default:
		return nil
	}
}

// GetUserByEmail for checking the existeence user
func (user User) GetUserByEmail(db *gorm.DB) (*User, error) {
	if err := db.Debug().Table("users").Preload("Role").Where("email = ?", user.Email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Register a new user
func (user *User) Register(db *gorm.DB) (*User, error) {
	var err error
	if err := db.Debug().Create(&user).Error; err != nil {
		return nil, err
	}
	return user, err
}

// GetUsers Get all users
func (userJSON UserJSON) GetUsers(db *gorm.DB) (*[]UserJSON, error) {
	var err error
	users := []UserJSON{}
	if err := db.Table("users").Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, err
}

// GetUser Get one user
func (userJSON UserJSON) GetUser(id string, db *gorm.DB) (*UserJSON, error) {
	if err := db.Debug().Table("users").Preload("Role").Where("id = ?", id).First(&userJSON).Error; err != nil {
		return nil, err
	}
	return &userJSON, nil
}

// ChangeUserPassword change user password
func (user *User) ChangeUserPassword(id string, db *gorm.DB) (*User, error) {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	if err := db.Debug().Table("users").Where("id = ?", id).Update("password", hashedPassword).Error; err != nil {
		return nil, err
	}
	return user, nil
}
