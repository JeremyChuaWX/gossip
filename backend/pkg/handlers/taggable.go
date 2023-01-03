package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TaggableHandler struct {
	DB *gorm.DB
}

func (h *TaggableHandler) CreateTaggable(c *fiber.Ctx) error {
	type tagPostInput struct {
		TagID  string `json:"tag_id" validate:"required"`
		PostID string `json:"post_id" validate:"required"`
	}

	var err error
	var input tagPostInput
	var post models.Post
	currId := utils.GetJwt(c)

	// bind input struct
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// get post by id
	if err = h.DB.Where("id = ?", input.PostID).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}

	// check authorised
	if currId != post.UserID {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
	}

	taggable := models.Taggable{
		TagID:  input.TagID,
		PostID: input.PostID,
	}

	// tag post
	if err = h.DB.Create(&taggable).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Post tagged",
		Data:  taggable,
	})
}
