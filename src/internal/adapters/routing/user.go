package routing

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kallabs/idp-api/src/internal/adapters/factory"
	"github.com/kallabs/idp-api/src/internal/utils"
)

type userHandler struct {
	BaseHanlder
}

func ConfigureUserHandler(db *sqlx.DB, r *mux.Router) {
	h := &userHandler{
		BaseHanlder: BaseHanlder{
			db:     db,
			router: r,
		},
	}

	h.router.HandleFunc("/me", h.Get()).Methods("GET")
	h.router.HandleFunc("/users/{username:[a-zA-Z0..9]+}", h.FindByUsername()).Methods("GET")
}

func (i *userHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := utils.LoggedInUser(r)
		if user == nil {
			log.Println("error user get context")
			return
		}

		utils.SendJson(w, user, http.StatusOK)
	}
}

func (i *userHandler) FindByUsername() http.HandlerFunc {
	get_by_username_service := factory.NewGetUserByUsernameService(i.db)

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		user := utils.LoggedInUser(r)
		if user == nil {
			utils.SendJsonError(w, "You don't have permissions.", http.StatusBadRequest)
			return
		}

		usr, err := get_by_username_service.Execute(username)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, usr, http.StatusOK)
	}
}
