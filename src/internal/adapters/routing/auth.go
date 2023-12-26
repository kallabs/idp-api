package routing

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kallabs/idp-api/src/internal/adapters/factory"
	value_objects "github.com/kallabs/idp-api/src/internal/domain"
	"github.com/kallabs/idp-api/src/internal/utils"
)

type authRouter struct {
	BaseHanlder
}

func ConfigureAuthHandler(db *sqlx.DB, r *mux.Router) {
	i := &authRouter{
		BaseHanlder: BaseHanlder{
			db:     db,
			router: r,
		},
	}

	i.router.HandleFunc("/signup", i.Signup()).Methods("POST")
	i.router.HandleFunc("/signup/{token}", i.ConfirmSignup()).Methods("POST")
	i.router.HandleFunc("/jwt/create", i.CreateJWT()).Methods("POST")
	i.router.HandleFunc("/jwt/refresh", i.RefreshJWT()).Methods("POST")
}

func (i *authRouter) Signup() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	create_reg_service := factory.NewCreateRegistrantService(i.db)

	return func(w http.ResponseWriter, r *http.Request) {
		var s request

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		email_addr := value_objects.EmailAddress(s.Email)

		if _, err := create_reg_service.Execute(email_addr, s.Username, s.Password); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		// TODO: send signup email

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *authRouter) ConfirmSignup() http.HandlerFunc {
	confirm_email_service := factory.NewConfirmEmailService(i.db)

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if err := confirm_email_service.Execute(vars["token"]); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *authRouter) CreateJWT() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	login_service := factory.NewLoginWithEmailPasswordService(i.db)

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
		}
		defer r.Body.Close()

		auth_payload, err := login_service.Execute(s.Email, s.Password)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		cookieAccessToken := http.Cookie{
			Name:     "_sso_access_token",
			Value:    auth_payload.AccessToken,
			Expires:  time.Now().Add(5 * time.Minute),
			Domain:   "localhost",
			HttpOnly: false,
		}
		cookieRefreshToken := http.Cookie{
			Name:     "_sso_refresh_token",
			Value:    auth_payload.RefreshToken,
			Expires:  time.Now().Add(24 * time.Hour),
			Domain:   "localhost",
			HttpOnly: true,
		}
		http.SetCookie(w, &cookieAccessToken)
		http.SetCookie(w, &cookieRefreshToken)

		utils.SendJson(w, auth_payload, http.StatusOK)
	}
}

func (i *authRouter) RefreshJWT() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	refresh_token_service := factory.NewRefreshTokensService()

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
		}
		defer r.Body.Close()

		auth_payload, err := refresh_token_service.Execute(s.RefreshToken)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, auth_payload, http.StatusOK)
	}
}
