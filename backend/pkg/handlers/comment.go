package handlers

import (
	"database/sql"
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/validate"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type createCommentInput struct {
	UserID   string `json:"user_id" validate:"required"`
	PostID   string `json:"post_id" validate:"required"`
	ParentID string `json:"parent_id,omitempty"`
	Body     string `json:"body" validate:"required"`
}

type updateCommentInput struct {
	Body string `json:"body"`
}

type CommentHandler struct {
	DB *gorm.DB
}

func (h CommentHandler) CreateComment(c *fiber.Ctx) error {
	var err error
	var input createCommentInput

	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := validate.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var parentId sql.NullString
	if input.ParentID != "" {
		parentId = sql.NullString{String: input.ParentID, Valid: true}
	} else {
		parentId = sql.NullString{Valid: false}
	}

	cmt := models.Comment{
		UserID:   input.UserID,
		PostID:   input.PostID,
		ParentID: parentId,
		Body:     input.Body,
	}

	// create comment
	if err = h.DB.Create(&cmt).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"data": cmt})
}

func (h CommentHandler) GetCommentById(c *fiber.Ctx) error {
	var err error
	var cmt models.Comment
	id := c.Params("id")

	// get comment by id
	if err = h.DB.Where("id = ?", id).Preload(clause.Associations).First(&cmt).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "Comment not found")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": cmt})
}

func (h CommentHandler) UpdateComment(c *fiber.Ctx) error {
	var err error
	var input updateCommentInput
	var cmt models.Comment
	id := c.Params("id")

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := validate.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// get comment by id
	if err = h.DB.Where("id = ?", id).First(&cmt).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "Comment not found")
	}

	updateCmt := models.Comment{Body: input.Body}

	// update comment
	if err = h.DB.Model(&cmt).Updates(updateCmt).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": cmt})
}

func (h CommentHandler) DeleteComment(c *fiber.Ctx) error {
	var err error
	var cmt models.Comment
	id := c.Params("id")

	// get comment by id
	if err = h.DB.Where("id = ?", id).First(&cmt).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "Comment not found")
	}

	// delete comment
	if err = h.DB.Delete(&cmt).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": true})
}
