package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type signUpInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required"`
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	DB *gorm.DB
}

func (h AuthHandler) SignUp(c *gin.Context) {
	var err error
	var input signUpInput

	// input validation
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    strings.ToLower(input.Email),
		Password: hashedPassword,
	}

	// create user
	if err = h.DB.Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h AuthHandler) SignIn(c *gin.Context) {
	var err error
	var input signInInput
	var user models.User

	// input validation
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	// get user by username
	if err = h.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// verify password
	if err = utils.VerifyPassword(user.Password, input.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	// update user session

	// update auth middleware

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h AuthHandler) SignOut(c *gin.Context) {}
