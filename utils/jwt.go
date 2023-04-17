package utils

import (
	"errors"
	"time"

	"readygo/pkg/settings"

	jwt "github.com/golang-jwt/jwt/v4"
)

// Claims Claims type
type Claims struct {
	Username    string   `json:"usr"`
	Permissions []string `json:"perms"`
	jwt.RegisteredClaims
}

// GenerateToken generate JWT token
func GenerateToken(username string, permissions []string) (string, error) {
	if settings.JWT.Secret == "" {
		return "", errors.New("JWT secret must be set")
	}
	nowTime := time.Now()
	expireTime := nowTime.Add(settings.JWT.Expires)

	claims := Claims{
		username,
		permissions,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    settings.JWT.Issuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(settings.JWT.Secret))

	return token, err
}

// ParseToken parse JWT token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, errors.New("JWT parsing error")
}
