package email

import (
	"github.com/kallabs/sso-api/src/internal/app/valueobject"
)

func (s *EmailService) SendSignup(email valueobject.EmailAddress, token string) error {
	subject := "Подтверждение регистрации"
	from := "admin@akarpovich.online"

	data := make(map[string]interface{})
	data["Token"] = token

	return s.SendWithView(
		subject,
		from,
		[]string{string(email)},
		[]string{
			"./assets/email/layout/base.html",
			"./assets/email/auth/signup.html",
		},
		"layout",
		data,
	)
}
