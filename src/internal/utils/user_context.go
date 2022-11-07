package utils

import (
	"net/http"

	"github.com/kallabs/sso-api/src/internal/app"
)

func LoggedInUser(r *http.Request) *app.User {
	user, ok := r.Context().Value("user").(*app.User)

	if !ok {
		return nil
	}

	return user
}
