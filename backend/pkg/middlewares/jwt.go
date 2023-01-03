package middlewares

import (
	"errors"
	"gossip/backend/pkg/config"
	"gossip/backend/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type jwtExtractor func(*fiber.Ctx) (string, error)

func jwtFromCookie(c *fiber.Ctx) (string, error) {
	rawToken := c.Cookies("access_token")

	if rawToken == "" {
		return "", errors.New("Missing or malformed JWT")
	}

	return rawToken, nil
}

func jwtFromHeader(c *fiber.Ctx) (string, error) {
	bearerToken := c.Get("Authorization")

	rawToken := strings.Split(bearerToken, " ")
	if len(rawToken) == 2 {
		return rawToken[1], nil
	}

	return "", errors.New("Missing or malformed JWT")
}

func getJwtExtractors() []jwtExtractor {
	return []jwtExtractor{
		jwtFromCookie,
		jwtFromHeader,
	}
}

func Jwtware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var rawToken string
		var err error

		for _, extractor := range getJwtExtractors() {
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

		sub, err := utils.ValidateToken(rawToken, env.AccessTokenPublicKey)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
		}

		c.Locals("user-id", sub)
		return c.Next()
	}
}
