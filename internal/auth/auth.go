package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CreateJwtToken(durationInSeconds int, userID string, tokenSecret string) (string, error) {
	deltaTime := time.Duration(durationInSeconds) * time.Second
	if deltaTime == 0 || deltaTime > 24*time.Hour {
		deltaTime = 24 * time.Hour
	}

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(deltaTime)),
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(tokenSecret))
	return signedString, err
}

func ValidateJwtToken(tokenString string, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return userIDString, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	authStrings := strings.Split(authHeader, " ")
	if len(authStrings) != 2 || authStrings[0] != "Bearer" {
		return "", errors.New("invalid auth header, expected: Bearer xxx")
	}

	return authStrings[1], nil
}
