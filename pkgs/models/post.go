package models

import (
	"errors"
	"html/template"

	"github.com/coderunne/jwt-login/pkgs/database"
	"gorm.io/gorm"
)

type Post struct {
	*gorm.Model
	Title      string        `json:"title" form:"title"`
	Slug       string        `json:"slug"`
	CategoryID uint          `json:"category_id" form:"category_id"`
	UserID     uint          `json:"user_id" form:"user_id"`
	Thumbnail  string        `json:"thumbnail" form:"thumbnail"`
	Tags       []*Tag        `json:"tags,omitempty" form:"tags" gorm:"many2many:post_tags"`
	Body       template.HTML `json:"body" form:"body"`
}

var (
	ErrCreatePost    = errors.New("Error Adding Post To The Database")
	ErrInvalidPostID = errors.New("ID provided is invalid!")
	ErrPostNotFound  = errors.New("Post was not found!")
)

func CreatePost(post *Post) error {
	result := database.DB.Model(&Post{}).Create(post)
	if result.RowsAffected < 1 {
		return ErrCreatePost
	}
	return nil
}

func GetAllPosts() ([]*Post, error) {
	var posts []*Post
	result := database.DB.Model(&Post{}).Preload("Tags").Find(&posts)
	if result.Error != nil {
		return posts, result.Error
	}
	return posts, nil
}

func GetPostByID(id uint) (*Post, error) {
	var post *Post
	if id < 1 {
		return &Post{}, ErrInvalidPostID
	}

	// Retrieve user using the provided id
	result := database.DB.Model(&Post{}).Preload("Tags").Where("id = ?", id).Scan(&post)
	if result.Error != nil {
		return &Post{}, result.Error
	} else if result.RowsAffected == 0 {
		return &Post{}, ErrPostNotFound
	}

	return post, nil
}

func GetPostBySlug(slug string) (*Post, error) {
	// Retrieve user using the provided id
	var post *Post
	result := database.DB.Model(&Post{}).Preload("Tags").Where("slug = ?", slug).Scan(&post)
	if result.Error != nil {
		return &Post{}, result.Error
	} else if result.RowsAffected == 0 {
		return &Post{}, ErrPostNotFound
	}
	return post, nil
}

func UpdatePost(slug string, post *Post) error {
	result := database.DB.Model(&Post{}).Where("slug = ?", slug).Updates(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeletePost(slug string) error {
	err := database.DB.Where("slug = ?", slug).Delete(&Post{}).Error
	if err != nil {
		return err
	}
	return nil
}
