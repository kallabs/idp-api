package routing

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type BaseHanlder struct {
	db     *sqlx.DB
	router *mux.Router
}

func (i *BaseHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.router.ServeHTTP(w, r)
}

func Index(w http.ResponseWriter, r *http.Request) {
	// Authenticate the request, get the id from the route params,
	// and fetch the user from the DB, etc.

	// Get the token and pass it in the CSRF header. Our JSON-speaking client
	// or JavaScript framework can now read the header and return the token in
	// in its own "X-CSRF-Token" request header on the subsequent POST.
	csrfToken := csrf.Token(r)
	w.Header().Set("X-CSRF-Token", csrfToken)
	w.Write([]byte(csrfToken))
}
