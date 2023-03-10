package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TagHandler struct {
	DB *gorm.DB
}

func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	type createTagInput struct {
		Name string `json:"name" validate:"required"`
	}

	var err error
	var input createTagInput

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

	tag := models.Tag{Name: input.Name}

	// create tag
	if err = h.DB.Create(&tag).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Tag created",
		Data:  tag,
	})
}

func (h *TagHandler) GetAllTags(c *fiber.Ctx) error {
	var err error
	var tags []models.Tag

	// create tag
	if err = h.DB.Find(&tags).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Tags found",
		Data:  tags,
	})
}

func (h *TagHandler) GetTagById(c *fiber.Ctx) error {
	var err error
	var tag models.Tag
	id := c.Params("id")

	// get tag by id
	if err = h.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Tag not found")
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Tag found",
		Data:  tag,
	})
}

func (h *TagHandler) UpdateTag(c *fiber.Ctx) error {
	type updateTagInput struct {
		Name string `json:"name,omitempty"`
	}

	var err error
	var tag models.Tag
	var input updateTagInput
	id := c.Params("id")

	// get tag by id
	if err = h.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Tag not found")
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

	updateTag := models.Tag{Name: input.Name}

	// update tag
	if err = h.DB.Model(&tag).Updates(updateTag).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Tag updated",
		Data:  tag,
	})
}

func (h *TagHandler) DeleteTag(c *fiber.Ctx) error {
	var err error
	var tag models.Tag
	id := c.Params("id")

	// get tag by id
	if err = h.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Tag not found")
	}

	// delete tag
	if err = h.DB.Delete(&tag).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Tag deleted",
		Data:  tag,
	})
}
