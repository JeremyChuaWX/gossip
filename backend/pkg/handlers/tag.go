package handlers

import (
	"gossip/backend/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createTagInput struct {
	Name string `binding:"required"`
}

type updateTagInput struct {
	Name string
}

type TagHandler struct {
	DB *gorm.DB
}

func (h TagHandler) CreateTag(c *gin.Context) {
	var err error
	var input createTagInput

	// input validation
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	tag := models.Tag{Name: input.Name}

	// create tag
	if err = h.DB.Create(&tag).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": tag})
}

func (h TagHandler) GetAllTags(c *gin.Context) {
	var err error
	var tags []models.Tag

	// create tag
	if err = h.DB.Find(&tags).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tags})
}

func (h TagHandler) GetTagById(c *gin.Context) {
	var err error
	var tag models.Tag
	id := c.Param("id")

	// get tag by id
	if err = h.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tag})
}

func (h TagHandler) UpdateTag(c *gin.Context) {
	var err error
	var tag models.Tag
	var input updateTagInput
	id := c.Param("id")

	// get tag by id
	if err = h.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// input validation
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	updateTag := models.Tag{Name: input.Name}

	// update tag
	if err = h.DB.Model(&tag).Updates(updateTag).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tag})
}

func (h TagHandler) DeleteTag(c *gin.Context) {
	var err error
	var tag models.Tag
	id := c.Param("id")

	// get tag by id
	if err = h.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// delete tag
	if err = h.DB.Delete(&tag).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
