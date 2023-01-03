package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	var err error
	var user models.User
	id := c.Params("id")
	currId := utils.GetJwt(c)

	// get user by id
	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
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

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	type updateUserInput struct {
		Email    string `json:"email,omitempty" validate:"omitempty,email"`
		Password string `json:"password,omitempty"`
	}

	var err error
	var input updateUserInput
	var user models.User
	id := c.Params("id")
	currId := utils.GetJwt(c)

	// get user by id
	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// check authorised
	if currId != user.ID {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
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

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	var err error
	var user models.User
	id := c.Params("id")
	currId := utils.GetJwt(c)

	// get user by id
	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// check authorised
	if currId != user.ID {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
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
	return fiber.NewError(fiber.StatusInternalServerError, "Not ready")
}
