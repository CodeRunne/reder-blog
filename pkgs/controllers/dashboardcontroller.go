package controllers

import (
	"fmt"
	"net/http"
	"strings"
	
	"github.com/coderunne/jwt-login/pkgs/models"
	"github.com/gin-gonic/gin"
)

func DashboardController(c *gin.Context) {

	// Get authenticated user info
	data, _ := c.Get("user")
	user := data.(*models.User)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Welcome to the dashboard! %s", strings.ToTitle(user.Fullname)),
	})
}
