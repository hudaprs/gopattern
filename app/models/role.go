package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Role Struct
type Role struct {
	gorm.Model
	Name string `gorm:"size:100;not null;"`
}

// GetRoles getting all roles
func (role *Role) GetRoles(begin, limit int, name string, db *gorm.DB) (*[]Role, error) {
	roles := []Role{}
	if err := db.Debug().Table("roles").Offset(begin).Limit(limit).Where("name LIKE ?", "%"+name+"%").Find(&roles).Error; err != nil {
		return nil, err
	}
	return &roles, nil
}

// Validate a input user
func (role Role) Validate() error {
	if role.Name == "" {
		return errors.New("Name is required")
	}
	return nil
}

// Create create a new role
func (role *Role) Create(db *gorm.DB) (*Role, error) {
	var err error
	if err := db.Debug().Table("roles").Create(&role).Error; err != nil {
		return nil, err
	}
	return role, err
}

// GetRoleByID get role by ID
func (role Role) GetRoleByID(id string, db *gorm.DB) (*Role, error) {
	if err := db.Debug().Table("roles").Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// Delete selected role
func (role *Role) Delete(id uint, db *gorm.DB) (*Role, error) {
	var err error
	if err := db.Debug().Table("roles").Where("id = ?", id).Unscoped().Delete(&role).Error; err != nil {
		return nil, err
	}
	return role, err
}

// CountRoles in database
func (role *Role) CountRoles(db *gorm.DB) (int, error) {
	var count int
	if err := db.Debug().Table("roles").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
