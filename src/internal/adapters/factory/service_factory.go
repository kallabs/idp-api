package factory

import (
	"github.com/jmoiron/sqlx"
	"github.com/kallabs/idp-api/src/internal/app/services"
	"github.com/kallabs/idp-api/src/internal/app/usecases"
)

func NewCreateRegistrantService(db *sqlx.DB) usecases.CreateRegistrant {
	return &services.CreateRegistrantService{
		UserGateway:     NewUserGateway(db),
		SendSignupEmail: NewSendSignupEmailService(),
	}
}

func NewConfirmEmailService(db *sqlx.DB) *services.ConfirmEmailService {
	return &services.ConfirmEmailService{
		UserGateway: NewUserGateway(db),
	}
}

func NewLoginWithEmailPasswordService(db *sqlx.DB) *services.LoginWithEmailPasswordService {
	return &services.LoginWithEmailPasswordService{
		UserGateway:           NewUserGateway(db),
		GenerateTokensService: NewGenerateTokensService(),
	}
}

func NewRefreshTokensService() *services.RefreshTokensService {
	return &services.RefreshTokensService{
		JwtGateway: DefaultJwtGateway(),
	}
}

func NewGetUserInfoService(db *sqlx.DB) *services.GetUserInfoService {
	return &services.GetUserInfoService{
		UserGateway: NewUserGateway(db),
	}
}

func NewGetUserByUsernameService(db *sqlx.DB) *services.GetUserByUsernameService {
	return &services.GetUserByUsernameService{
		UserGateway: NewUserGateway(db),
	}
}

func NewSendSignupEmailService() *services.SendSignupEmailService {
	return &services.SendSignupEmailService{
		EmailGateway: NewEmailGateway(),
	}
}

func NewGenerateTokensService() *services.GenerateTokensService {
	return &services.GenerateTokensService{
		JwtGateway: DefaultJwtGateway(),
	}
}
