package middlewares

import (
	"gossip/backend/pkg/config"
	"gossip/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func Jwtware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var rawToken string
		var err error

		for _, extractor := range utils.GetJwtExtractors() {
			rawToken, err = extractor(c)

			// break when jwt found
			if rawToken != "" && err == nil {
				break
			}
		}

		env, err := config.GetEnv()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Cannot get environment variables")
		}

		sub, err := utils.ValidateJwt(rawToken, env.AccessTokenPublicKey)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
		}

		if err = utils.SetUserId(c, sub); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error setting user id")
		}

		return c.Next()
	}
}
