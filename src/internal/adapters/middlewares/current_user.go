package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kallabs/idp-api/src/internal/adapters/factory"
	"github.com/kallabs/idp-api/src/internal/utils"
)

func CurrentUser(db *sqlx.DB) func(http.Handler) http.Handler {
	user_gateway := factory.NewUserGateway(db)
	jwt_gateway := factory.DefaultJwtGateway()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := r.Header.Get("Authorization")

			if len(token) > 0 {
				splitToken := strings.Split(token, "Bearer ")
				token = splitToken[1]

				claims, err := jwt_gateway.ParseTokenPayload(token)

				if err != nil {
					utils.SendJsonError(w, "invalid token", http.StatusUnauthorized)
					return
				}

				if claims != nil {
					user, err := user_gateway.GetById(claims.UserId)

					if err == nil {
						ctx := context.WithValue(r.Context(), "user", user)

						r = r.WithContext(ctx)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
