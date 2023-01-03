package handlers

import (
	"gossip/backend/pkg/config"
	"gossip/backend/pkg/models"
	"gossip/backend/pkg/utils"
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
	if errors := utils.ValidateStruct(&input); errors != nil {
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

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "User created",
		Data:  user,
	})
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
	if errors := utils.ValidateStruct(&input); errors != nil {
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

	// get env
	env, err := config.GetEnv()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot get environment variables")
	}

	// generate access token
	accessToken, err := utils.CreateJwt(env.AccessTokenDuration, user.ID, env.AccessTokenPrivateKey)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// generate refresh token
	refreshToken, err := utils.CreateJwt(env.RefreshTokenDuration, user.ID, env.RefreshTokenPrivateKey)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// set cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   env.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   env.RefreshTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   env.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: false,
	})

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "User signed in",
		Data:  user,
	})
}

func (h *AuthHandler) RefreshAccessToken(c *fiber.Ctx) error {
	var err error
	var user models.User

	cookie := c.Cookies("refresh_token")
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, "Could not refresh access token")
	}

	// load env
	env, err := config.GetEnv()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot get environment variables")
	}

	sub, err := utils.ValidateJwt(cookie, env.RefreshTokenPublicKey)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	// get user by username
	if err = h.DB.Where("username = ?", sub).First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	accessToken, err := utils.CreateJwt(env.AccessTokenDuration, user.ID, env.AccessTokenPrivateKey)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   env.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   env.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: false,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": accessToken})
}

func (h *AuthHandler) SignOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   -1,
		Secure:   false,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   -1,
		Secure:   false,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   -1,
		Secure:   false,
		HTTPOnly: false,
	})

	return c.Status(fiber.StatusOK).JSON(models.ServerResponse{
		Error: false,
		Msg:   "User signed out",
	})
}
