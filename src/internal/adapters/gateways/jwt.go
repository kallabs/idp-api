package gateways

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	value_objects "github.com/kallabs/idp-api/src/internal/domain"
	"github.com/kallabs/idp-api/src/internal/domain/entities"
	"github.com/kallabs/idp-api/src/internal/utils"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	Uid  value_objects.ID   `json:"uid"`
	Type entities.TokenType `json:"type"`
}

func newToken(
	userId value_objects.ID,
	rsa_private_key string,
	claims *TokenClaims,
) string {
	private_key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(rsa_private_key))
	method := jwt.GetSigningMethod(utils.Conf.Jwt.Algorithm)
	token := jwt.NewWithClaims(method, claims)

	tokenString, err := token.SignedString(private_key)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(private_key)

	return tokenString
}

type JwtGateway struct {
	AccessLifetime  int
	RefreshLifetime int
	Issuer          string
	Audience        []string
	RsaPrivateKey   string
	RsaPublicKey    string
}

func (i *JwtGateway) NewAccessToken(user_id value_objects.ID) string {
	claims := &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(i.AccessLifetime) * time.Minute)),
			Issuer:    i.Issuer,
			Audience:  jwt.ClaimStrings(i.Audience),
		},
		user_id,
		entities.TypeAccess,
	}

	return newToken(user_id, i.RsaPrivateKey, claims)
}

func (i *JwtGateway) NewRefreshToken(user_id value_objects.ID) string {
	claims := &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(i.RefreshLifetime) * time.Minute)),
			Issuer:    i.Issuer,
			Audience:  jwt.ClaimStrings(i.Audience),
		},
		user_id,
		entities.TypeRefresh,
	}

	return newToken(user_id, i.RsaPrivateKey, claims)
}

func (i *JwtGateway) ParseTokenPayload(token_string string) (*entities.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(token_string, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM([]byte(i.RsaPublicKey))
	})

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return &entities.TokenPayload{
			UserId: claims.Uid,
			Type:   claims.Type,
		}, nil
	}

	return nil, err
}
