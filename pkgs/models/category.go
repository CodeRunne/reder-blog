package models

import (
	"errors"

	"github.com/coderunne/jwt-login/pkgs/database"
	"gorm.io/gorm"
)

type Category struct {
	*gorm.Model
	Name  string `json:"name"`
	Posts []Post `json:"posts,omitempty"`
}

var (
	ErrInvalidCategoryID = errors.New("Provided Categrory ID does not exist")
	ErrCategoryResultNotFound     = errors.New("Category ID not found in the database")
	ErrCreateCategory    = errors.New("Error Adding Category To The Database")
)

func CreateCategory(category *Category) error {
	result := database.DB.Model(&Category{}).Create(category)
	if result.RowsAffected < 1 {
		return ErrCreateCategory
	}
	return nil
}

func GetCategory(id uint) (*Category, error) {
	var category *Category
	if id < 1 {
		return &Category{}, ErrInvalidCategoryID
	}

	// Retrieve user using the provided id
	result := database.DB.Model(&Category{}).Where("id = ?", id).Scan(&category)
	if result.Error != nil {
		return &Category{}, result.Error
	} else if result.RowsAffected == 0 {
		return &Category{}, ErrCategoryResultNotFound
	}

	return category, nil
}

func GetAllCategory() ([]*Category, error) {
	var category []*Category
	result := database.DB.Model(&Category{}).Select("id, name, created_at").Find(&category)
	if result.Error != nil {
		return category, result.Error
	}
	return category, nil
}

func UpdateCategory(id uint, name string) error {
	result := database.DB.Model(&Category{}).Where("id = ?", id).Updates(&Category{Name: name})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetCategoryPosts(id uint) ([]Post, error) {
	var category *Category
	err := database.DB.Model(&Category{}).Where("id = ?", id).Preload("Posts").Find(&category).Error
	if err != nil {
		return category.Posts, err
	}
	return category.Posts, nil
}

func DeleteCategory(id uint) error {
	err := database.DB.Delete(&Category{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
