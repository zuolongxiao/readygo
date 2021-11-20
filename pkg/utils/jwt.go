package utils

import (
	"errors"
	"time"

	"readygo/pkg/settings"

	jwt "github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(settings.AppSetting.JwtSecret)
var jwtExpires = settings.AppSetting.JwtExpires
var jwtIssuer = settings.AppSetting.JwtIssuer

// Claims Claims type
type Claims struct {
	Username    string   `json:"usr"`
	Permissions []string `json:"perms"`
	jwt.RegisteredClaims
}

// GenerateToken generate JWT token
func GenerateToken(username string, permissions []string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(jwtExpires * time.Hour)

	claims := Claims{
		username,
		permissions,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    jwtIssuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken parse JWT token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, errors.New("JWT parsing error")
}
