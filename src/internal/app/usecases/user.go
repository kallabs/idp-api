package usecases

import (
	value_objects "github.com/kallabs/idp-api/src/internal/domain"
	"github.com/kallabs/idp-api/src/internal/domain/entities"
)

type GetUserInfo interface {
	Execute(user_id value_objects.ID) (*entities.UserInfo, error)
}

type GetUserByUsername interface {
	Execute(string) (*entities.UserInfo, error)
}
