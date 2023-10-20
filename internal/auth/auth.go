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

func CreateJwtToken(tokenType string, userID string, tokenSecret string) (string, error) {
	issuer := "chirpy"
	deltaTime := time.Hour
	switch typ := tokenType; typ {
	case "access":
		issuer += "-access"
	case "refresh":
		issuer += "-refresh"
		deltaTime = time.Hour * 24 * 60
	default:
		return "", errors.New("invalid token type")
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(deltaTime)),
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(tokenSecret))
	return signedString, err
}

func ValidateJwtToken(tokenString string, tokenSecret string, issuer string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	tokenIssuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}

	if tokenIssuer != issuer {
		return "", errors.New("invalid token issuer")
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
