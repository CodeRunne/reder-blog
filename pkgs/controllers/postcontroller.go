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
	cID, _ := strconv.Atoi(c.PostForm("category_id"))
	body := c.PostForm("body")

	var tags []models.Tag
	if strings.Contains(c.PostForm("tags"), ",") {
		tags = append(tags, strings.Split(c.PostForm("tags"), ",")...)
	} else {
		tags = append(tags, strings.Fields(strings)...)
	}

	fmt.Println(tags)

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
	fmt.Println(filename, ext)
	path := filepath.Join("storage/posts", filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return
	}

	// check if category exists
	_, err = models.GetCategory(uint(cID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Assign data to struct
	post.Title = title
	post.Slug = strings.ToLower(strings.Replace(title, " ", "-", -1))
	post.CategoryID = uint(cID)
	post.UserID = user.ID
	post.Body = template.HTML(body)
	// post.Tags = append(post.Tags, tags)
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
	cID, _ := strconv.Atoi(c.PostForm("category_id"))
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
	_, err = models.GetCategory(uint(cID))
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
	post.CategoryID = uint(cID)
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
