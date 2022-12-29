package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func FiberError(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	if err = c.Status(code).JSON(e); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	return nil
}
