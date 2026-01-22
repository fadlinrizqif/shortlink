package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hashedPass, nil
}

func CheckPassword(password, hash string) (bool, error) {
	isPassCorrect, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return isPassCorrect, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	clams := &jwt.RegisteredClaims{
		Issuer:    "shortlink",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clams)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return ss, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("Invalid Token")
	}

	subjectId, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, errors.New("Something wrong in claims")
	}

	validID, err := uuid.Parse(subjectId)
	if err != nil {
		return uuid.Nil, err
	}

	return validID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	rawtoken := headers.Get("Authorization")
	if rawtoken == "" {
		return "", errors.New("Token not found")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(rawtoken, prefix) {
		return "", errors.New("invalid authorization header")
	}

	sanitizeToken := strings.Replace(rawtoken, prefix, "", 1)

	return sanitizeToken, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}
