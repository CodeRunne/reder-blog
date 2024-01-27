package controllers

import (
	"net/http"
	"strconv"

	"github.com/coderunne/jwt-login/pkgs/models"
	"github.com/gin-gonic/gin"
)

func CategoriesCreateController(c *gin.Context) {
	var category models.Category

	// Bind data to struct
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := models.CreateCategory(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category Created Successfully",
	})
}

func CategoriesAllController(c *gin.Context) {
	categories, err := models.GetAllCategory()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func CategoriesUpdateController(c *gin.Context) {
	var category models.Category
	// Get url param
	id, _ := strconv.Atoi(c.Param("id"))
	// bind json data to category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update category
	err := models.UpdateCategory(uint(id), category.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category has successfully updated",
	})
}

func CategoriesDeleteController(c *gin.Context) {
	// Get url param
	id, _ := strconv.Atoi(c.Param("id"))
	if err := models.DeleteCategory(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.String(http.StatusOK, "Category deleted successfully")
}

func CategoryPostsController(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	posts, err := models.GetCategoryPosts(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, posts)
}
