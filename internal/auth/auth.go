package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSigned, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenSigned, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(subject)

}

func GetBearerToken(headers http.Header) (string, error) {

	authString := headers.Get("Authorization")
	if authString == "" {
		return "", errors.New("no authorization present")
	}
	authWords := strings.Fields(authString)
	if len(authWords) != 2 || authWords[0] != "Bearer" {
		return "", errors.New("malformed authorization string")
	}
	token_string := authWords[1]

	return token_string, nil

}

func MakeRefreshToken() (string, error) {
	nums := make([]byte, 32)
	_, err := rand.Read(nums)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(nums)
	return token, nil
}

func GetAPIKey(headers http.Header) (string, error) {

	authString := headers.Get("Authorization")
	if authString == "" {
		return "", errors.New("no authorization present")
	}
	authWords := strings.Fields(authString)
	if len(authWords) != 2 || authWords[0] != "ApiKey" {
		return "", errors.New("malformed authorization string")
	}
	apiKey := authWords[1]

	return apiKey, nil

}
