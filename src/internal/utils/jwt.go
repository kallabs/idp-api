package utils

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kallabs/sso-api/src/internal/app/valueobject"
)

type TokenType uint8

const (
	TypeAccess TokenType = iota
	TypeRefresh
)

type TokenClaims struct {
	jwt.RegisteredClaims
	Uid  *valueobject.ID `json:"uid"`
	Type TokenType       `json:"type"`
}

func newToken(userId *valueobject.ID, claims *TokenClaims) string {
	hmacSecret := []byte(Conf.SecretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(hmacSecret)

	return tokenString
}

func NewAccessToken(userId *valueobject.ID) string {
	claims := &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(Conf.Jwt.AccessLifetime) * time.Minute)),
			Issuer:    Conf.Jwt.Issuer,
			Audience:  jwt.ClaimStrings(Conf.Jwt.Audience),
		},
		userId,
		TypeAccess,
	}

	return newToken(userId, claims)
}

func NewRefreshToken(userId *valueobject.ID) string {
	claims := &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(Conf.Jwt.RefreshLifetime) * time.Minute)),
			Issuer:    Conf.Jwt.Issuer,
			Audience:  jwt.ClaimStrings(Conf.Jwt.Audience),
		},
		userId,
		TypeRefresh,
	}

	return newToken(userId, claims)
}

func GetTokenClaims(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(Conf.SecretKey), nil
	})

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
