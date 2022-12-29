package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"
	"gossip/backend/pkg/validate"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type updateUserInput struct {
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty"`
}

type UserHandler struct {
	DB *gorm.DB
}

func (h UserHandler) GetUserById(c *fiber.Ctx) error {
	var err error
	var user models.User
	id := c.Params("id")

	// get user by id
	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "User not found")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": user})
}

func (h UserHandler) UpdateUser(c *fiber.Ctx) error {
	var err error
	var input updateUserInput
	var user models.User
	id := c.Params("id")

	// get user by id
	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "User not found")
	}

	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := validate.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// only hash new password if provided
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	updateUser := models.User{
		Email:    input.Email,
		Password: hashedPassword,
	}

	// update user
	if err = h.DB.Model(&user).Updates(updateUser).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": user})
}

func (h UserHandler) DeleteUser(c *fiber.Ctx) error {
	var err error
	var user models.User
	id := c.Params("id")

	// get user by id
	if err = h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return fiber.NewError(http.StatusNotFound, "User not found")
	}

	// delete user
	if err = h.DB.Delete(&user).Error; err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": true})
}

func (h UserHandler) ToggleProfileVisibility(c *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusInternalServerError, "Not ready")
}
