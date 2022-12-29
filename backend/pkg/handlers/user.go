package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type updateUserInput struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password"`
}

type UserHandler struct {
	DB *gorm.DB
}

func (h UserHandler) GetUserById(c *gin.Context) {
	var err error
	var user models.User
	id := c.Param("id")

	// get user by id
	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h UserHandler) UpdateUser(c *gin.Context) {
	var err error
	var input updateUserInput
	var user models.User
	id := c.Param("id")

	// get user by id
	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// input validation
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	hashedPassword := ""

	// only hash new password if provided
	if input.Password != "" {
		hashedPassword, err = utils.HashPassword(input.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	updateUser := models.User{
		Email:    input.Email,
		Password: hashedPassword,
	}

	// update user
	if err = h.DB.Model(&user).Updates(updateUser).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h UserHandler) DeleteUser(c *gin.Context) {
	var err error
	var user models.User
	id := c.Param("id")

	// get user by id
	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// delete user
	if err = h.DB.Delete(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (h UserHandler) ToggleProfileVisibility(c *gin.Context) {
}
