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
	cookie := c.Cookies("access_token")

	if cookie == "" {
		return "", errors.New("Missing or malformed JWT")
	}

	return cookie, nil
}

func jwtFromHeader(c *fiber.Ctx) (string, error) {
	bearerToken := c.Get("Authorization")
	l := len("Bearer")

	if len(bearerToken) > l+1 && strings.EqualFold(bearerToken[:l], "Bearer") {
		return strings.TrimSpace(bearerToken[l:]), nil
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
		var auth string
		var err error

		for _, extractor := range getJwtExtractors() {
			auth, err = extractor(c)

			// break when jwt found
			if auth != "" && err == nil {
				break
			}
		}

		env, err := config.GetEnv()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Cannot get environment variables")
		}

		sub, err := utils.ValidateToken(auth, env.AccessTokenPublicKey)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorised")
		}

		c.Set("user_id", sub)
		return c.Next()
	}
}
