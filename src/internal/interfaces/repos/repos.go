package repos

import (
	"github.com/kallabs/idp-api/src/internal/app"
	"github.com/kallabs/idp-api/src/internal/interfaces/db"
)

type Repos struct {
	User app.UserRepo
}

func NewRepos(db db.DB) *Repos {
	return &Repos{
		User: NewUserRepo(db),
	}
}
