package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type jwtExtractor func(*fiber.Ctx) (string, error)

func CreateJwt(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	var err error

	// decode base64 privateKey
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("Could not decode key: %w", err)
	}

	// parse rsa private key
	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("Could not parse key from PEM: %w", err)
	}

	now := time.Now().UTC()

	// create claims
	claims := jwt.MapClaims{
		"sub": payload,
		"exp": now.Add(ttl).Unix(),
	}

	// sign signed
	signed, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(rsaKey)
	if err != nil {
		return "", fmt.Errorf("Could not sign token: %w", err)
	}

	return signed, nil
}

func ValidateJwt(token string, publicKey string) (string, error) {
	var err error

	// throw if empty token
	if token == "" {
		return "", fmt.Errorf("Missing token")
	}

	// decode base64 privateKey
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("Could not decode key: %w", err)
	}

	// parse rsa public key
	rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return "", fmt.Errorf("Could not parse key from PEM: %w", err)
	}

	// function that returns the public key used to sign the token
	var keyFunc jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return rsaKey, nil
	}

	// parse token string to jwt
	parsedToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return "", err
	}

	// cast claims struct
	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	// throw if invalid
	if !ok || !parsedToken.Valid || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return "", fmt.Errorf("Invalid token")
	}

	return claims["sub"].(string), nil
}

func jwtFromCookie(c *fiber.Ctx) (string, error) {
	rawToken := c.Cookies("access_token")

	if rawToken == "" {
		return "", fmt.Errorf("Missing or malformed JWT")
	}

	return rawToken, nil
}

func jwtFromHeader(c *fiber.Ctx) (string, error) {
	bearerToken := c.Get("Authorization")

	rawToken := strings.Split(bearerToken, " ")
	if len(rawToken) == 2 {
		return rawToken[1], nil
	}

	return "", fmt.Errorf("Missing or malformed JWT")
}

func GetJwtExtractors() []jwtExtractor {
	return []jwtExtractor{
		jwtFromCookie,
		jwtFromHeader,
	}
}

func SetUserId(c *fiber.Ctx, sub string) error {
	if sub == "" {
		return fmt.Errorf("Empty subject")
	}

	c.Locals("user-id", sub)

	return nil
}

func GetUserId(c *fiber.Ctx) string {
	id := c.Locals("user-id")

	// handle not signed in
	if id != nil {
		return c.Locals("user-id").(string)
	}

	return ""
}
