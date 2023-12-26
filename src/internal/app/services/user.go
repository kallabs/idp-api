package services

import (
	"log"

	"github.com/kallabs/idp-api/src/internal/app/usecases"
	value_objects "github.com/kallabs/idp-api/src/internal/domain"
	"github.com/kallabs/idp-api/src/internal/domain/entities"
)

type GetUserInfoService struct {
	UserGateway usecases.GetUserInfoGateway
}

func (i *GetUserInfoService) Execute(user_id value_objects.ID) (*entities.UserInfo, error) {
	user, err := i.UserGateway.GetById(user_id)
	if err != nil {
		log.Println(err)
	}

	user_info := &entities.UserInfo{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return user_info, nil
}

type GetUserByUsernameService struct {
	UserGateway usecases.GetUserByUsernameGateway
}

func (i *GetUserByUsernameService) Execute(username string) (*entities.UserInfo, error) {
	user, err := i.UserGateway.GetByUsername(username)
	if err != nil {
		log.Println(err)
	}

	user_info := &entities.UserInfo{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return user_info, nil
}
