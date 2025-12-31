package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

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

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}

	authHeaderParts := strings.SplitN(authHeader, " ", 2)
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return "", http.ErrNoCookie
	}

	return authHeaderParts[1], nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	
	return hex.EncodeToString(key) 
}