package utils

import (
	"net/http"

	"github.com/kallabs/idp-api/src/internal/domain/entities"
)

func LoggedInUser(r *http.Request) *entities.User {
	user, ok := r.Context().Value("user").(*entities.User)

	if !ok {
		return nil
	}

	return user
}
