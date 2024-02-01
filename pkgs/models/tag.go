package models

import (
	"errors"

	"github.com/coderunne/jwt-login/pkgs/database"
	"gorm.io/gorm"
)

type Tag struct {
	*gorm.Model
	Name  string `json:"name"`
	Posts []Post `json:"posts,omitempty" gorm:"many2many:post_tags"`
}

var (
	ErrInvalidTagID      = errors.New("Provided Tag ID does not exist")
	ErrTagResultNotFound = errors.New("Tag ID not found in the database")
	ErrCreateTag         = errors.New("Error Adding Tag To The Database")
)

func CreateTag(tag *Tag) error {
	result := database.DB.Model(&Tag{}).Create(tag)
	if result.RowsAffected < 1 {
		return ErrCreateTag
	}
	return nil
}

func GetTag(id uint) (*Tag, error) {
	var tag *Tag
	if id < 1 {
		return &Tag{}, ErrInvalidTagID
	}

	// Retrieve user using the provided id
	result := database.DB.Model(&Tag{}).Where("id = ?", id).Scan(&tag)
	if result.Error != nil {
		return &Tag{}, result.Error
	} else if result.RowsAffected == 0 {
		return &Tag{}, ErrTagResultNotFound
	}

	return tag, nil
}

func GetAllTags() ([]*Tag, error) {
	var tags []*Tag
	result := database.DB.Model(&Tag{}).Select("id, name, created_at").Find(&tags)
	if result.Error != nil {
		return tags, result.Error
	}
	return tags, nil
}

func UpdateTag(id uint, name string) error {
	result := database.DB.Model(Tag{}).Where("id = ?", id).Updates(Tag{Name: name})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetTagPosts(id uint) ([]Post, error) {
	var tag *Tag
	err := database.DB.Model(&Tag{}).Where("id = ?", id).Preload("Posts").Find(&tag).Error
	if err != nil {
		return tag.Posts, err
	}
	return tag.Posts, nil
}

func DeleteTag(id uint) error {
	err := database.DB.Delete(&Tag{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
