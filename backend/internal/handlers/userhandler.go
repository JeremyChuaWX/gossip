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
	var err error
	var user models.User
	id := c.Param("id")

	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h Handler) CreateUser(c *gin.Context) {
	var err error
	var input createUserInput

	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user := models.User{Username: input.Username, Email: input.Email, Password: input.Password}

	if err = h.DB.Create(&user).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h Handler) UpdateUser(c *gin.Context) {
	var err error
	var user models.User
	id := c.Param("id")

	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var input updateUserInput

	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	h.DB.Model(&user).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h Handler) DeleteUser(c *gin.Context) {
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
