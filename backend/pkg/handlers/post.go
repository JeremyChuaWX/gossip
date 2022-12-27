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
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (h PostHandler) CreatePost(c *gin.Context) {
	var err error
	var input createPostInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	post := models.Post{
		UserID: input.UserID,
		Title:  input.Title,
		Body:   input.Body,
	}

	if err = h.DB.Create(&post).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": post})
}

func (h PostHandler) GetAllPosts(c *gin.Context) {
	var err error
	var posts []models.Post

	if err = h.DB.Find(&posts).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Posts not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func (h PostHandler) GetPost(c *gin.Context) {
	var err error
	var post models.Post
	id := c.Param("id")

	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Post not found"})
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
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	updatePost := models.Post{
		Title: input.Title,
		Body:  input.Body,
	}

	if err = h.DB.Model(&post).Updates(updatePost).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": post})
}

func (h PostHandler) DeletePost(c *gin.Context) {
	var err error
	var post models.Post
	id := c.Param("id")

	if err = h.DB.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err = h.DB.Delete(&post).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
