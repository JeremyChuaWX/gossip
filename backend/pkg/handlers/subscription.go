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

	subscription := models.Subscription{
		UserID: currId,
		PostID: post.ID,
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
	var subscription models.Subscription

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

	// get subscription by ids
	if err = h.DB.Where("user_id = ? AND post_id = ?", currId, input.PostID).First(&subscription).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Taggable not found")
	}

	// check authorised
	if currId != subscription.UserID {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
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
