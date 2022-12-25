package handlers

import (
	"gossip/backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}

type updateUserInput struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}

func (h Handler) GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h Handler) CreateUser(c *gin.Context) {
	var input createUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Username: input.Username, Email: input.Email, Password: input.Password}
	h.DB.Create(&user)

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h Handler) UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	var input updateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Model(&user).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h Handler) DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	h.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
