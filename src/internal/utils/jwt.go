package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kallabs/sso-api/src/internal/app/valueobject"
)

var hmacSecret []byte = []byte(os.Getenv("SECRET_KEY"))

type TokenClaims struct {
	jwt.RegisteredClaims
	Uid *valueobject.ID `json:"uid"`
}

func NewToken(userId *valueobject.ID) string {
	claims := TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			Issuer:    "sso.kallabs.by",
			Audience:  jwt.ClaimStrings{"kallabs.by", "lst.kallabs.by"},
		},
		userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}

func GetTokenClaims(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return hmacSecret, nil
	})

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
