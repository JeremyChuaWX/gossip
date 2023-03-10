package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserHandler struct {
	DB *gorm.DB
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	var err error
	var user models.User
	id := c.Params("id")

	// get user id (allow empty for not signed in)
	currId := utils.GetUserId(c)

	// get user by id
	if err = h.DB.Where("id = ?", id).Preload(clause.Associations).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// check authorised
	if !user.IsPublic && currId != user.ID {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "User found",
		Data:  user,
	})
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	var err error
	var user models.User

	// get user id
	currId := utils.GetUserId(c)
	if currId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	// get user by id
	if err = h.DB.Where("id = ?", currId).Preload(clause.Associations).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "User found",
		Data:  user,
	})
}

func (h *UserHandler) UpdateMe(c *fiber.Ctx) error {
	type updateUserInput struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty" validate:"omitempty,email"`
		Password string `json:"password,omitempty"`
	}

	var err error
	var input updateUserInput
	var user models.User

	// get user id
	currId := utils.GetUserId(c)
	if currId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	// get user by id
	if err = h.DB.Where("id = ?", currId).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// bind input struct
	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := utils.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// only hash new password if provided
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	updateUser := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	// update user
	if err = h.DB.Model(&user).Updates(updateUser).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "User updated",
		Data:  user,
	})
}

func (h *UserHandler) DeleteMe(c *fiber.Ctx) error {
	var err error
	var user models.User

	// get user id
	currId := utils.GetUserId(c)
	if currId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	// get user by id
	if err = h.DB.Where("id = ?", currId).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// delete user
	if err = h.DB.Delete(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "User deleted",
		Data:  user,
	})
}

func (h *UserHandler) ToggleProfileVisibility(c *fiber.Ctx) error {
	var err error
	var user models.User

	// get user id
	currId := utils.GetUserId(c)
	if currId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	// get user by id
	if err = h.DB.Where("id = ?", currId).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	updateUser := models.User{
		IsPublic: !user.IsPublic,
	}

	// toggle visibility
	if err = h.DB.Model(&user).Select("IsPublic").Updates(updateUser).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "User profile visibility toggled",
		Data:  user,
	})
}
