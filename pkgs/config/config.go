package config

import (
	"log"
	"os"

	"github.com/coderunne/jwt-login/pkgs/database"
	"github.com/coderunne/jwt-login/pkgs/models"
	"github.com/coderunne/jwt-login/pkgs/utility"
	"github.com/gin-gonic/gin"
)

func init() {

	// Set gin default mode to release
	gin.SetMode(gin.ReleaseMode)

	if err := utility.Load("./.env"); err != nil {
		log.Fatal(err.Error())
	}

	var (
		DB_NAME = os.Getenv("DB_NAME")
		DB_USER = os.Getenv("DB_USER")
		DB_PASS = os.Getenv("DB_PASS")
		DB_HOST = os.Getenv("DB_HOST")
	)

	if err := database.Connect(DB_USER, DB_PASS, DB_HOST, DB_NAME); err != nil {
		log.Fatal(err.Error())
	}

	database.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Category{}, &models.Tag{}, &models.Post{})

}
