package controllers

import (
	"net/http"
	"strconv"

	"github.com/coderunne/jwt-login/pkgs/models"
	"github.com/gin-gonic/gin"
)

func TagCreateController(c *gin.Context) {
	var tag models.Tag

	// Bind data to struct
	if err := c.BindJSON(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := models.CreateTag(&tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tag Created Successfully",
	})
}

func TagsAllController(c *gin.Context) {
	tags, err := models.GetAllTags()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func TagUpdateController(c *gin.Context) {
	var tag models.Tag
	// Get url param
	id, _ := strconv.Atoi(c.Param("id"))
	// bind json data to tag
	if err := c.BindJSON(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update tag
	err := models.UpdateTag(uint(id), tag.Name)
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

func TagDeleteController(c *gin.Context) {
	// Get url param
	id, _ := strconv.Atoi(c.Param("id"))
	if err := models.DeleteTag(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.String(http.StatusOK, "Tag deleted successfully")
}

func TagPostsController(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	posts, err := models.GetTagPosts(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, posts)
}
