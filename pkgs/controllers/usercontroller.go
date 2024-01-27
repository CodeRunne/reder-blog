package controllers

import (
	"net/http"
	"strconv"

	"github.com/coderunne/jwt-login/pkgs/models"
	"github.com/gin-gonic/gin"
)

func UsersController(c *gin.Context) {
	users, _ := models.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

func UserGetController(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := models.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, user)
}

func UserPostsController(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	posts, err := models.GetUserPosts(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, posts)
}
