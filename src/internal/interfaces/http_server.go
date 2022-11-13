package interfaces

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/kallabs/sso-api/src/internal/app/usecases"
	app_handlers "github.com/kallabs/sso-api/src/internal/interfaces/handlers"
	"github.com/kallabs/sso-api/src/internal/interfaces/middlewares"
	"github.com/kallabs/sso-api/src/internal/interfaces/repos"
	"github.com/kallabs/sso-api/src/internal/interfaces/services"
)

func configureRouter(repos *repos.Repos, services *services.Services) http.Handler {
	baseRouter := mux.NewRouter().StrictSlash(true)
	protectedRouter := baseRouter.PathPrefix("/in").Subrouter()

	csrfMiddleware := csrf.Protect(
		[]byte(os.Getenv("SECRET_KEY")),
		// set Secure true only for production
		csrf.Secure(false),
		// use * only for development purpose
		csrf.TrustedOrigins([]string{"*"}),
		// instruct the browser to never send cookies during cross site requests
		// csrf.SameSite(csrf.SameSiteStrictMode),
	)

	protectedRouter.Use(csrfMiddleware)

	baseRouter.HandleFunc("/csrf", app_handlers.Index).Methods("GET")

	authInterector := usecases.NewAuthInteractor(repos.User, services.Email)
	app_handlers.ConfigureAuthHandler(authInterector, baseRouter)

	userInterector := usecases.NewUserInteractor(repos.User)
	app_handlers.ConfigureUserHandler(userInterector, protectedRouter)

	return baseRouter
}

// NewServer - initialize HTTP Server
func NewHTTPServer(address string, repos *repos.Repos, services *services.Services) (*http.Server, error) {
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
		middlewares.CurrentUser(repos.User),
		//mw.Authorizer(authEnforcer),
	).Then(configureRouter(repos, services))

	server := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      chain,
		Addr:         address,
	}

	return server, nil
}
