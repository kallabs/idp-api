package services

import (
	"github.com/kallabs/sso-api/src/internal/app/services"
	"github.com/kallabs/sso-api/src/internal/interfaces/repos"
	"github.com/kallabs/sso-api/src/internal/interfaces/services/email"
)

type Services struct {
	Email services.EmailService
}

func NewServices(repos *repos.Repos) *Services {
	return &Services{
		Email: &email.EmailService{},
	}
}
