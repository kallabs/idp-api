package repos

import (
	"github.com/kallabs/sso-api/src/internal/app"
	"github.com/kallabs/sso-api/src/internal/interfaces/db"
)

type Repos struct {
	User app.UserRepo
}

func NewRepos(db db.DB) *Repos {
	return &Repos{
		User: NewUserRepo(db),
	}
}
