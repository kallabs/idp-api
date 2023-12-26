package usecases

import (
	value_objects "github.com/kallabs/idp-api/src/internal/domain"
	"github.com/kallabs/idp-api/src/internal/domain/entities"
)

type CreateUserGateway interface {
	Create(*entities.CreateUserSchema) (*entities.User, error)
}

type ConfirmEmailGateway interface {
	GetByRegToken(string) (*entities.User, error)
	UpdateFromEntity(user *entities.User) error
}

type LoginWithEmailPasswordGateway interface {
	GetByEmail(value_objects.EmailAddress) (*entities.User, error)
}

type GetUserInfoGateway interface {
	GetById(value_objects.ID) (*entities.User, error)
}

type GetUserByUsernameGateway interface {
	GetByUsername(string) (*entities.User, error)
}
