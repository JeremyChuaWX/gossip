package handlers

import (
	"gossip/backend/pkg/auth"
	"gossip/backend/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type signUpInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required"`
}

type updateUserInput struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password"`
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserHandler struct {
	DB *gorm.DB
}

func (h UserHandler) SignUp(c *gin.Context) {
	var err error
	var input signUpInput

	// input validation
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	// create user record in DB
	if err = h.DB.Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h UserHandler) GetUserById(c *gin.Context) {
	var err error
	var user models.User
	id := c.Param("id")

	// get user record from DB (by id)
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

	// get user record from DB (by id)
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
		hashedPassword, err = auth.HashPassword(input.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	updateUser := models.User{
		Email:    input.Email,
		Password: hashedPassword,
	}

	// update user record in DB
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

	// get user record from DB (by id)
	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// delete user record in DB
	if err = h.DB.Delete(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (h UserHandler) SignIn(c *gin.Context) {
	var err error
	var input signInInput
	var user models.User

	// input validation
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	// get user record from DB (by username)
	if err = h.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// verify password
	if err = auth.VerifyPassword(user.Password, input.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	// update user session

	// update auth middleware
}
