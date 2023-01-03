package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
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

func ValidateToken(token string, publicKey string) (string, error) {
	var err error

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
	if !ok || !parsedToken.Valid || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return "", fmt.Errorf("Invalid token")
	}

	return claims["sub"].(string), nil
}
