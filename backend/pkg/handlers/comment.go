package handlers

import (
	"database/sql"
	"gossip/backend/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type createCommentInput struct {
	UserID   int    `json:"user_id" binding:"required"`
	PostID   int    `json:"post_id" binding:"required"`
	ParentID int    `json:"parent_id"`
	Body     string `json:"body" binding:"required"`
}

type updateCommentInput struct {
	Body string `json:"body"`
}

type CommentHandler struct {
	DB *gorm.DB
}

func (h CommentHandler) CreateComment(c *gin.Context) {
	var err error
	var input createCommentInput

	// input validation
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	// provide null ParentID if none is provided (default is 0)
	var cmt models.Comment
	if input.ParentID != 0 {
		cmt = models.Comment{
			UserID:   input.UserID,
			PostID:   input.PostID,
			ParentID: sql.NullInt32{Int32: int32(input.ParentID), Valid: true},
			Body:     input.Body,
		}
	} else {
		cmt = models.Comment{
			UserID:   input.UserID,
			PostID:   input.PostID,
			ParentID: sql.NullInt32{Valid: false},
			Body:     input.Body,
		}
	}

	// create comment
	if err = h.DB.Create(&cmt).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": cmt})
}

func (h CommentHandler) GetCommentById(c *gin.Context) {
	var err error
	var cmt models.Comment
	id := c.Param("id")

	// get comment by id
	if err = h.DB.Where("id = ?", id).Preload(clause.Associations).First(&cmt).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cmt})
}

func (h CommentHandler) UpdateComment(c *gin.Context) {
	var err error
	var input updateCommentInput
	var cmt models.Comment
	id := c.Param("id")

	// input validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid fields"})
		return
	}

	// get comment by id
	if err = h.DB.Where("id = ?", id).First(&cmt).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	updateCmt := models.Comment{Body: input.Body}

	// update comment
	if err = h.DB.Model(&cmt).Updates(updateCmt).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cmt})
}

func (h CommentHandler) DeleteComment(c *gin.Context) {
	var err error
	var cmt models.Comment
	id := c.Param("id")

	// get comment by id
	if err = h.DB.Where("id = ?", id).First(&cmt).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// delete comment
	if err = h.DB.Delete(&cmt).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
