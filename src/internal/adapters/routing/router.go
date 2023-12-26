package routing

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kallabs/idp-api/src/internal/utils"
)

func ConfigureRouter(db *sqlx.DB) http.Handler {
	baseRouter := mux.NewRouter().StrictSlash(true)
	protectedRouter := baseRouter.PathPrefix("/in").Subrouter()

	csrfMiddleware := csrf.Protect(
		[]byte(utils.Conf.SecretKey),
		// set Secure true only for production
		csrf.Secure(false),
		// use * only for development purpose
		csrf.TrustedOrigins([]string{"*"}),
		// instruct the browser to never send cookies during cross site requests
		// csrf.SameSite(csrf.SameSiteStrictMode),
	)

	protectedRouter.Use(csrfMiddleware)

	baseRouter.HandleFunc("/csrf", Index).Methods("GET")
	ConfigureAuthHandler(db, baseRouter)
	ConfigureUserHandler(db, protectedRouter)

	return baseRouter
}
