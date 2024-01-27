package models

import (
	"github.com/coderunne/jwt-login/pkgs/database"
	"gorm.io/gorm"
)

type Role struct {
	*gorm.Model
	Name   string `json:"name,omitempty"`
	UserID uint   `json:"user_id,omitempty" gorm:"unique"`
}

func CreateRole(id uint, name string) error {
	result := database.DB.Model(&Role{}).Create(&Role{
		Name:   name,
		UserID: id,
	})

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateRole(id uint, name string) error {
	result := database.DB.Model(&Role{}).Where("user_id = ?", id).Updates(&Role{
		Name: name,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
