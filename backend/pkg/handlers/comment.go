package handlers

import (
	"database/sql"
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentHandler struct {
	DB *gorm.DB
}

func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	type createCommentInput struct {
		PostID   string `json:"post_id" validate:"required"`
		ParentID string `json:"parent_id,omitempty"`
		Body     string `json:"body" validate:"required"`
	}

	var err error
	var input createCommentInput
	currId := utils.GetJwt(c)

	// bind input struct
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var parentId sql.NullString
	if input.ParentID != "" {
		parentId = sql.NullString{String: input.ParentID, Valid: true}
	} else {
		parentId = sql.NullString{Valid: false}
	}

	cmt := models.Comment{
		UserID:   currId,
		PostID:   input.PostID,
		ParentID: parentId,
		Body:     input.Body,
	}

	// create comment
	if err = h.DB.Create(&cmt).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": cmt})
}

func (h *CommentHandler) GetCommentById(c *fiber.Ctx) error {
	var err error
	var cmt models.Comment
	id := c.Params("id")

	// get comment by id
	if err = h.DB.Where("id = ?", id).Preload(clause.Associations).First(&cmt).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Comment not found")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": cmt})
}

func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	type updateCommentInput struct {
		Body         string `json:"body,omitempty"`
		CommentScore int    `json:"comment_score,omitempty"`
	}

	var err error
	var input updateCommentInput
	var cmt models.Comment
	id := c.Params("id")

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// get comment by id
	if err = h.DB.Where("id = ?", id).First(&cmt).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Comment not found")
	}

	updateCmt := models.Comment{
		Body:         input.Body,
		CommentScore: input.CommentScore,
	}

	// update comment
	if err = h.DB.Model(&cmt).Updates(updateCmt).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": cmt})
}

func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	var err error
	var cmt models.Comment
	id := c.Params("id")

	// get comment by id
	if err = h.DB.Where("id = ?", id).First(&cmt).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Comment not found")
	}

	// delete comment
	if err = h.DB.Delete(&cmt).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": true})
}
