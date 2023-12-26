package factory

import (
	"github.com/jmoiron/sqlx"
	"github.com/kallabs/idp-api/src/internal/adapters/gateways"
	"github.com/kallabs/idp-api/src/internal/adapters/gateways/postgres"
	"github.com/kallabs/idp-api/src/internal/utils"
)

func NewUserGateway(db *sqlx.DB) *postgres.UserGateway {
	return &postgres.UserGateway{
		Db: db,
	}
}

func NewEmailGateway() *gateways.EmailGateway {
	return &gateways.EmailGateway{}
}

func NewJwtGateway(
	access_lifetime int,
	refresh_lifetime int,
	issuer string,
	audience []string,
	rsa_private_key string,
	rsa_public_key string,
) *gateways.JwtGateway {
	return &gateways.JwtGateway{
		AccessLifetime:  access_lifetime,
		RefreshLifetime: refresh_lifetime,
		Issuer:          issuer,
		Audience:        audience,
		RsaPrivateKey:   rsa_private_key,
		RsaPublicKey:    rsa_public_key,
	}
}

func DefaultJwtGateway() *gateways.JwtGateway {
	return NewJwtGateway(
		utils.Conf.Jwt.AccessLifetime,
		utils.Conf.Jwt.RefreshLifetime,
		utils.Conf.Jwt.Issuer,
		utils.Conf.Jwt.Audience,
		utils.Conf.RsaPrivateKey,
		utils.Conf.RsaPublicKey,
	)
}
