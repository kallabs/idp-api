package services

import (
	"errors"
	"log"
	"time"

	"github.com/kallabs/idp-api/src/internal/app/usecases"
	value_objects "github.com/kallabs/idp-api/src/internal/domain"
	"github.com/kallabs/idp-api/src/internal/domain/entities"
	"github.com/kallabs/idp-api/src/internal/utils"
	"github.com/kallabs/idp-api/src/pkg"
)

type CreateRegistrantService struct {
	UserGateway     usecases.CreateUserGateway
	SendSignupEmail usecases.SendSignupEmail
}

func (i *CreateRegistrantService) Execute(email value_objects.EmailAddress, username string, password string) (*entities.User, error) {
	inRegistrant := entities.CreateUserSchema{
		Email:          value_objects.EmailAddress(email),
		Username:       username,
		Token:          pkg.RandomString(utils.Conf.RegTokenLength),
		Status:         entities.UserUnconfirmed,
		TokenExpiresAt: time.Now().Add(12 * time.Hour).UTC(),
	}
	inRegistrant.SetPassword(password)

	registrant, err := i.UserGateway.Create(&inRegistrant)
	if err != nil {
		return nil, err
	}

	go i.SendSignupEmail.Execute(registrant.Email, registrant.Token)

	return registrant, nil
}

type ConfirmEmailService struct {
	UserGateway usecases.ConfirmEmailGateway
}

func (i *ConfirmEmailService) Execute(token string) error {
	registrant, err := i.UserGateway.GetByRegToken(token)
	if err != nil {
		log.Println(err)
		return err
	}

	registrant.Status = entities.UserActive

	if err := i.UserGateway.UpdateFromEntity(registrant); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

type GenerateTokensService struct {
	JwtGateway usecases.NewTokensGateway
}

func (i *GenerateTokensService) Execute(user_id value_objects.ID) (*entities.AuthPayload, error) {
	auth_payload := &entities.AuthPayload{
		AccessToken:  i.JwtGateway.NewAccessToken(user_id),
		RefreshToken: i.JwtGateway.NewRefreshToken(user_id),
	}
	return auth_payload, nil
}

type RefreshTokensService struct {
	JwtGateway usecases.RefreshTokensGateway
}

func (i *RefreshTokensService) Execute(refresh_token string) (*entities.AuthPayload, error) {

	token_payload, err := i.JwtGateway.ParseTokenPayload(refresh_token)
	if err != nil {
		return nil, err
	}

	if token_payload.Type != entities.TypeRefresh {
		return nil, errors.New("Invalid token type")
	}

	auth_payload := &entities.AuthPayload{
		AccessToken:  i.JwtGateway.NewAccessToken(token_payload.UserId),
		RefreshToken: i.JwtGateway.NewRefreshToken(token_payload.UserId),
	}

	return auth_payload, nil
}

type LoginWithEmailPasswordService struct {
	UserGateway           usecases.LoginWithEmailPasswordGateway
	GenerateTokensService usecases.GenerateTokens
}

func (i *LoginWithEmailPasswordService) Execute(email string, password string) (*entities.AuthPayload, error) {
	email_addr := value_objects.EmailAddress(email)
	user, err := i.UserGateway.GetByEmail(email_addr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("incorrect password")
	}

	return i.GenerateTokensService.Execute(*user.Id)
}
