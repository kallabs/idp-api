package usecases

import (
	value_objects "github.com/kallabs/idp-api/src/internal/domain"
	"github.com/kallabs/idp-api/src/internal/domain/entities"
)

type NewTokensGateway interface {
	NewAccessToken(user_id value_objects.ID) string
	NewRefreshToken(user_id value_objects.ID) string
}

type RefreshTokensGateway interface {
	ParseTokenPayload(string) (*entities.TokenPayload, error)
	NewAccessToken(user_id value_objects.ID) string
	NewRefreshToken(user_id value_objects.ID) string
}

type CreateRegistrant interface {
	Execute(email value_objects.EmailAddress, username string, password string) (*entities.User, error)
}

type ConfirmEmail interface {
	Execute(token string) error
}

type GenerateTokens interface {
	Execute(user_id value_objects.ID) (*entities.AuthPayload, error)
}

type RefreshTokens interface {
	Execute(refresh_token string) (*entities.AuthPayload, error)
}

type LoginWithEmailPassword interface {
	Execute(email string, password string) (*entities.AuthPayload, error)
}
