package services

import (
	"github.com/kallabs/idp-api/src/internal/app/usecases"
	value_objects "github.com/kallabs/idp-api/src/internal/domain"
)

type SendSignupEmailService struct {
	EmailGateway usecases.SendEmailWithViewGateway
}

func (i *SendSignupEmailService) Execute(email value_objects.EmailAddress, token string) error {
	subject := "Подтверждение регистрации"
	from := "admin@akarpovich.online"

	data := make(map[string]interface{})
	data["Token"] = token

	return i.EmailGateway.SendWithView(
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
