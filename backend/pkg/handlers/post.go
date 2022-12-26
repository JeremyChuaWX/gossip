package handlers

import (
	"gossip/backend/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostHandler struct {
	DB *gorm.DB
}

type createPostInput struct {
	UserID int    `json:"user_id" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
}

type updatePostInput struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

func (h PostHandler) CreatePost(c *gin.Context) {
	var err error
	var input createPostInput

	c.BindJSON(&input)

	post := models.Post{
		UserID: input.UserID,
		Title:  input.Title,
		Body:   input.Body,
	}

	if err = h.DB.Create(&post).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": post})
}

func (h PostHandler) GetAllPosts(c *gin.Context) {
	var err error
	var posts []models.Post

	if err = h.DB.Find(posts).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func (h PostHandler) GetPost(c *gin.Context) {
	var err error
	var post models.Post
	id := c.Param("id")

	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": post})
}

func (h PostHandler) UpdatePost(c *gin.Context) {
	var err error
	var input updatePostInput
	var post models.Post
	id := c.Param("id")

	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.BindJSON(&input)

	h.DB.Model(&post).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": post})
}

func (h PostHandler) DeletePost(c *gin.Context) {
	var err error
	var post models.Post
	id := c.Param("id")

	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err = h.DB.Delete(&post).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
