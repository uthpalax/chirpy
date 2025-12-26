package auth

import (
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func ComparePasswordAndHash(password, hash string) (bool, error) {
	compaire, err := argon2id.ComparePasswordAndHash(password, hash)

	return compaire, err
}

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims {
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userId.String(),
	})

	tokenString, err := token.SignedString([]byte(tokenSecret))
	return tokenString, err
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, err
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}
