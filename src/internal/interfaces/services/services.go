package services

import (
	"github.com/kallabs/idp-api/src/internal/app/services"
	"github.com/kallabs/idp-api/src/internal/interfaces/repos"
	"github.com/kallabs/idp-api/src/internal/interfaces/services/email"
)

type Services struct {
	Email services.EmailService
}

func NewServices(repos *repos.Repos) *Services {
	return &Services{
		Email: &email.EmailService{},
	}
}
