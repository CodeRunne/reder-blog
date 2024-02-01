package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/coderunne/jwt-login/pkgs/models"
	"github.com/gin-gonic/gin"
)

func PostAllController(c *gin.Context) {
	posts, err := models.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func PostGetController(c *gin.Context) {
	slug := c.Param("slug")
	post, err := models.GetPostBySlug(slug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, post)
}

func PostCreateController(c *gin.Context) {
	var post models.Post

	// Get authenticated user info
	data, _ := c.Get("user")
	user := data.(*models.User)

	// Get Form Data
	title := c.PostForm("title")
	categoryID, _ := strconv.Atoi(c.PostForm("category_id"))
	body := c.PostForm("body")
	post_tags := c.PostForm("tags")

	// Tags array
	var tags []*models.Tag

	// convert string to array
	ids := strings.Split(strings.TrimSpace(post_tags), ",")
	for _, id := range ids {
		// Convert to int
		tag_id, _ := strconv.Atoi(id)
		// check if tag id exists
		tag, err := models.GetTag(uint(tag_id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// append to tags array
		tags = append(tags, tag)
	}

	// Get post thumbnail file
	file, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Store post thumbnail
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join("storage/posts", filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	// check if category exists
	_, err = models.GetCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Assign data to struct
	post.Title = title
	post.Slug = strings.ToLower(strings.Replace(title, " ", "-", -1))
	post.CategoryID = uint(categoryID)
	post.UserID = user.ID
	post.Body = template.HTML(body)
	post.Tags = tags
	post.Thumbnail = filename

	// Create Post
	err = models.CreatePost(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Post created Successfully",
	})

}

func PostUpdateController(c *gin.Context) {
	// Get url param
	slug := c.Param("slug")
	// Get Form Data
	title := c.PostForm("title")
	categoryID, _ := strconv.Atoi(c.PostForm("category_id"))
	body := c.PostForm("body")

	// Get post thumbnail file
	file, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	gpost, err := models.GetPostBySlug(slug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get working directory
	dir, _ := os.Getwd()
	// Get Filepath
	oldpath := filepath.Join(dir, "storage/posts", gpost.Thumbnail)
	// Remove image from filepath
	if err := os.Remove(oldpath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Store post thumbnail
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join("storage/posts", filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	// check if category exists
	_, err = models.GetCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Assign data to struct
	var post models.Post
	post.Title = title
	post.Slug = strings.ToLower(strings.Replace(title, " ", "-", -1))
	post.CategoryID = uint(categoryID)
	post.Body = template.HTML(body)
	post.Thumbnail = filename

	if err := models.UpdatePost(slug, &post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.String(http.StatusOK, "Post updated successfully!")

}

func PostDeleteController(c *gin.Context) {
	// Get url param
	slug := c.Param("slug")
	if err := models.DeletePost(slug); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.String(http.StatusOK, "Post deleted successfully")
}
