package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/validate"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type createTagInput struct {
	Name string `json:"name" validate:"required"`
}

type updateTagInput struct {
	Name string `json:"name,omitempty"`
}

type TagHandler struct {
	DB *gorm.DB
}

func (h TagHandler) CreateTag(c *fiber.Ctx) error {
	var err error
	var input createTagInput

	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := validate.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	tag := models.Tag{Name: input.Name}

	// create tag
	if err = h.DB.Create(&tag).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"data": tag})
}

func (h TagHandler) GetAllTags(c *fiber.Ctx) error {
	var err error
	var tags []models.Tag

	// create tag
	if err = h.DB.Find(&tags).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": tags})
}

func (h TagHandler) GetTagById(c *fiber.Ctx) error {
	var err error
	var tag models.Tag
	id := c.Params("id")

	// get tag by id
	if err = h.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "Tag not found")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": tag})
}

func (h TagHandler) UpdateTag(c *fiber.Ctx) error {
	var err error
	var tag models.Tag
	var input updateTagInput
	id := c.Params("id")

	// get tag by id
	if err = h.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "Tag not found")
	}

	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := validate.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	updateTag := models.Tag{Name: input.Name}

	// update tag
	if err = h.DB.Model(&tag).Updates(updateTag).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": tag})
}

func (h TagHandler) DeleteTag(c *fiber.Ctx) error {
	var err error
	var tag models.Tag
	id := c.Params("id")

	// get tag by id
	if err = h.DB.Where("id = ?", id).First(&tag).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "Tag not found")
	}

	// delete tag
	if err = h.DB.Delete(&tag).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": true})
}
