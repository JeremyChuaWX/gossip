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
	type createTaggableInput struct {
		TagID  string `json:"tag_id" validate:"required"`
		PostID string `json:"post_id" validate:"required"`
	}

	var err error
	var input createTaggableInput
	var post models.Post

	// get user id
	currId := utils.GetUserId(c)
	if currId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	// bind input struct
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ServerResponse{
			Error: true,
			Msg:   "Invalid input",
			Data:  errors,
		})
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

	// create taggable
	if err = h.DB.Create(&taggable).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Create taggable",
		Data:  taggable,
	})
}

func (h *TaggableHandler) DeleteTaggable(c *fiber.Ctx) error {
	type deleteTaggableInput struct {
		TagID  string `json:"tag_id" validate:"required"`
		PostID string `json:"post_id" validate:"required"`
	}

	var err error
	var input deleteTaggableInput
	var post models.Post
	var taggable models.Taggable

	// get user id
	currId := utils.GetUserId(c)
	if currId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	// bind input struct
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ServerResponse{
			Error: true,
			Msg:   "Invalid input",
			Data:  errors,
		})
	}

	// get post by id
	if err = h.DB.Where("id = ?", input.PostID).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}

	// check authorised
	if currId != post.UserID {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
	}

	// get taggable by ids
	if err = h.DB.Where("post_id = ? AND tag_id = ?", input.PostID, input.TagID).First(&taggable).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Taggable not found")
	}

	// delete taggable
	if err = h.DB.Delete(&taggable).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Delete taggable",
		Data:  taggable,
	})
}
