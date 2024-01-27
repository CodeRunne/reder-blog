package models

import (
	"errors"
	"sync"

	"github.com/coderunne/jwt-login/pkgs/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	Fullname string `json:"fullname" form:"fullname"`
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email" gorm:"unique"`
	Role     Role   `json:"role,omitempty" form:"role" gorm:"default:viewers"`
	Posts    []Post `json:"posts,omitempty" form:"posts"`
	Profile  string `json:"profile", form:"profile"`
	Password string `json:"password,omitempty" form:"password"`
}

var (
	ErrEmptyFields   = errors.New("User Fields Cannot Be Empty!")
	ErrInvalidUserID = errors.New("ID provided is invalid!")
)

var RwMutex sync.RWMutex

func Sanitize(c *gin.Context, u *User) error {
	RwMutex.Lock()
	defer RwMutex.Unlock()
	if u.Fullname == "" || u.Username == "" || u.Email == "" || u.Profile == "" || u.Password == "" {
		return ErrEmptyFields
	}
	return nil
}

func CreateUser(user *User) error {
	result := database.DB.Model(&User{}).Create(user)
	if result.RowsAffected < 1 {
		return result.Error
	}

	if user.ID == 1 {
		AssignRole(user.ID, "admin")
	}

	return nil
}

func GetAllUsers() ([]*User, error) {
	var users []*User

	// Retrieve all users from the database
	result := database.DB.Model(&User{}).Preload("Role").Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	return users, nil
}

func GetUserByID(id uint) (*User, error) {
	var user *User
	if id < 1 {
		return &User{}, ErrInvalidUserID
	}

	// Retrieve user using the provided id
	result := database.DB.Model(User{}).Where("id = ?", id).Scan(&user)
	if result.Error != nil {
		return &User{}, result.Error
	}

	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user *User

	// Retrieve user using the provided email
	result := database.DB.Model(&User{}).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return &User{}, result.Error
	}

	return user, nil
}

func GetUserPosts(id uint) ([]Post, error) {
	var user *User
	err := database.DB.Model(&User{}).Where("id = ?", id).Preload("Posts").Find(&user).Error
	if err != nil {
		return user.Posts, err
	}
	return user.Posts, nil
}

func AssignRole(id uint, role string) error {
	err := CreateRole(id, role)
	if err != nil {
		return err
	}
	return nil
}
