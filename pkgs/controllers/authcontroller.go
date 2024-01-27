package controllers

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/coderunne/jwt-login/pkgs/models"
	"github.com/coderunne/jwt-login/pkgs/utility"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordNotMatch = errors.New("Provided password does not match!")
	ErrProfileNotFound = errors.New("Profile picture is required for account!")
	ErrBinding          = errors.New("Error binding data!")
)

func RegisterController(c *gin.Context) {
	var user models.User

	// Get form data
	fullname := c.PostForm("fullname")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Source
	file, err := c.FormFile("profile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrProfileNotFound,
		})
		return
	}

	// Store user profile
	filename := filepath.Base(file.Filename)
	path := filepath.Join("storage/profiles", filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	// Bind data to user struct
	user.Fullname = fullname
	user.Username = username
	user.Email = email
	user.Profile = path
	user.Password = password

	// Sanitize string
	if err := models.Sanitize(c, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate Email
	if err := utility.ValidateEmail(user.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Encrypt password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Store hashed password in user struct
	user.Password = string(hash)

	// Add user to the database
	err = models.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send Mail to client on successfull registration
	err = utility.SendMail(user.Email, "Registration Successfull", "Click link to login to dashboard <a href='http://localhost:8000/login'>http://localhost:8000/login</a>")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User added successfully",
	})
}

func LoginController(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate Email
	if err := utility.ValidateEmail(user.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate email exists
	db_user, err := models.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Compare passwords
	if err = bcrypt.CompareHashAndPassword([]byte(db_user.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrPasswordNotMatch,
		})
		return
	}

	// Create JWT token
	token, err := utility.CreateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set Token Header
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24*30, "", "", false, true)

	// Redirect user to dashboard
	c.String(http.StatusOK, "You are logged in!")
}
