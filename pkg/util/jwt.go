package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/utf6/goApi/pkg/config"
	"time"
)

var jwtSecret = []byte(config.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	expireTime := time.Now().Add(3 * time.Hour).Unix()

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Username:       username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer: "goApi",
		},
	})
	token, err := tokenClaims.SignedString(jwtSecret)

	result["token"] = token
	result["expire"] = expireTime
	return result, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}