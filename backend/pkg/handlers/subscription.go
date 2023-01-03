package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	DB *gorm.DB
}

func (h *SubscriptionHandler) CreateSubscription(c *fiber.Ctx) error {
	type createSubscriptionInput struct {
		PostID string `json:"post_id" validate:"required"`
	}

	var err error
	var input createSubscriptionInput
	var post models.Post
	// currId := utils.GetJwt(c)

	// bind input struct
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// get post by id
	if err = h.DB.Where("id = ?", post.ID).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}

	subscription := models.Subscription{
		UserID: post.ID,
		PostID: input.PostID,
	}

	// create subscription
	if err = h.DB.Create(&subscription).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Create subscription",
		Data:  subscription,
	})
}

func (h *SubscriptionHandler) DeleteSubscription(c *fiber.Ctx) error {
	type deleteSubscriptionInput struct {
		PostID string `json:"post_id" validate:"required"`
	}

	var err error
	var input deleteSubscriptionInput
	var post models.Post
	// currId := utils.GetJwt(c)

	// bind input struct
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// get post by id
	if err = h.DB.Where("id = ?", post.ID).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Post not found")
	}

	subscription := models.Subscription{
		UserID: post.ID,
		PostID: input.PostID,
	}

	// delete subscription
	if err = h.DB.Delete(&subscription).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "Delete subscription",
		Data:  subscription,
	})
}
