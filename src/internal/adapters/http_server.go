package adapters

import (
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
	"github.com/kallabs/idp-api/src/internal/adapters/middlewares"
	"github.com/kallabs/idp-api/src/internal/adapters/routing"
)

// NewServer - initialize HTTP Server
func NewHTTPServer(address string, db *sqlx.DB) (*http.Server, error) {
	// authEnforcer, err := casbin.NewEnforcer(
	// 	"config/auth_model.conf",
	// 	"config/policy.csv")

	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"})

	chain := alice.New(
		handlers.CORS(headersOk, originsOk, methodsOk),
		middlewares.CurrentUser(db),
		//mw.Authorizer(authEnforcer),
	).Then(routing.ConfigureRouter(db))

	server := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      chain,
		Addr:         address,
	}

	return server, nil
}
