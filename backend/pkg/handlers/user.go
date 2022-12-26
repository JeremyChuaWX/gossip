package handlers

import (
	"gossip/backend/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

type createUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required"`
}

type updateUserInput struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password"`
}

func (h UserHandler) CreateUser(c *gin.Context) {
	var err error
	var input createUserInput

	c.BindJSON(&input)

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if err = h.DB.Create(&user).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h UserHandler) GetUser(c *gin.Context) {
	var err error
	var user models.User
	id := c.Param("id")

	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h UserHandler) UpdateUser(c *gin.Context) {
	var err error
	var input updateUserInput
	var user models.User
	id := c.Param("id")

	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.BindJSON(&input)

	h.DB.Model(&user).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h UserHandler) DeleteUser(c *gin.Context) {
	var err error
	var user models.User
	id := c.Param("id")

	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err = h.DB.Delete(&user).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
