package usecases

import (
	"log"

	"github.com/kallabs/sso-api/src/internal/app"
	"github.com/kallabs/sso-api/src/internal/app/valueobject"
)

type UserInteractor struct {
	UserRepo app.UserRepo
}

func NewUserInteractor(ur app.UserRepo) *UserInteractor {
	return &UserInteractor{ur}
}

func (i *UserInteractor) Get(userId *valueobject.ID) (*app.User, error) {
	user, err := i.UserRepo.Get(userId)
	if err != nil {
		log.Println(err)
	}

	return user, nil
}

func (i *UserInteractor) FindByUsername(username string) (*app.User, error) {
	user, err := i.UserRepo.FindByUsername(username)
	if err != nil {
		log.Println(err)
	}

	return user, nil
}
