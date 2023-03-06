package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kallabs/sso-api/src/internal/app"
	"github.com/kallabs/sso-api/src/internal/app/usecases"
	"github.com/kallabs/sso-api/src/internal/utils"
)

type AuthInteractor interface {
	CreateRegistrant(*usecases.Registrant) (*app.User, error)
	ConfirmEmail(token string) error
	FindUserByEmailPassword(email string, password string) (*app.User, error)
}

type authHanlder struct {
	BaseHanlder
	authInteractor AuthInteractor
}

func ConfigureAuthHandler(ai AuthInteractor, r *mux.Router) {
	i := &authHanlder{
		BaseHanlder: BaseHanlder{
			router: r,
		},
		authInteractor: ai,
	}

	i.router.HandleFunc("/signup", i.Signup()).Methods("POST")
	i.router.HandleFunc("/signup/{token}", i.ConfirmSignup()).Methods("POST")
	i.router.HandleFunc("/jwt/create", i.CreateJWT()).Methods("POST")
	i.router.HandleFunc("/jwt/refresh", i.RefreshJWT()).Methods("POST")
}

func (i *authHanlder) Signup() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request

		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, "Invalid request data", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		registrant := &usecases.Registrant{
			Username: s.Username,
			Email:    s.Email,
			Password: s.Password,
		}

		if _, err := i.authInteractor.CreateRegistrant(registrant); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		// TODO: send signup email

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *authHanlder) ConfirmSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if err := i.authInteractor.ConfirmEmail(vars["token"]); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		utils.SendJson(w, "Success", http.StatusOK)
	}
}

func (i *authHanlder) CreateJWT() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
		}
		defer r.Body.Close()

		user, err := i.authInteractor.FindUserByEmailPassword(s.Email, s.Password)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		accessToken := utils.NewAccessToken(user.Id)
		refreshToken := utils.NewRefreshToken(user.Id)

		resData := map[string]string{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		}

		cookieAccessToken := http.Cookie{
			Name:     "_sso_access_token",
			Value:    accessToken,
			Expires:  time.Now().Add(5 * time.Minute),
			Domain:   "localhost",
			HttpOnly: false,
		}
		cookieRefreshToken := http.Cookie{
			Name:     "_sso_refresh_token",
			Value:    refreshToken,
			Expires:  time.Now().Add(24 * time.Hour),
			Domain:   "localhost",
			HttpOnly: true,
		}
		http.SetCookie(w, &cookieAccessToken)
		http.SetCookie(w, &cookieRefreshToken)

		utils.SendJson(w, resData, http.StatusOK)
	}
}

func (i *authHanlder) RefreshJWT() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refreshToken"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var s request
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
		}
		defer r.Body.Close()

		refreshClaims, err := utils.GetTokenClaims(s.RefreshToken)
		if err != nil {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		if refreshClaims.Type != utils.TypeRefresh {
			utils.SendJsonError(w, err, http.StatusBadRequest)
			return
		}

		accessToken := utils.NewAccessToken(refreshClaims.Uid)
		refreshToken := utils.NewRefreshToken(refreshClaims.Uid)

		utils.SendJson(w, map[string]string{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		}, http.StatusOK)
	}
}
