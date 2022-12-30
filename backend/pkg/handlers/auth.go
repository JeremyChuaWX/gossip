package handlers

import (
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"
	"gossip/backend/pkg/validate"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	type signUpInput struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email,omitempty" validate:"omitempty,email"`
		Password string `json:"password" validate:"required"`
	}

	var err error
	var input signUpInput

	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := validate.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user := models.User{
		Username: input.Username,
		Email:    strings.ToLower(input.Email),
		Password: hashedPassword,
	}

	// create user
	if err = h.DB.Create(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": user})
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	type signInInput struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var err error
	var input signInInput
	var user models.User

	if err = c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid fields")
	}

	// input validation
	if errors := validate.ValidateStruct(&input); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// get user by username
	if err = h.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// verify password
	if err = utils.VerifyPassword(user.Password, input.Password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid username or password")
	}

	// update user session

	// update auth middleware

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": user})
}

func (h AuthHandler) SignOut(c *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusInternalServerError, "Not ready")
}
