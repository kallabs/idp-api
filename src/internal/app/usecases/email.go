package usecases

import value_objects "github.com/kallabs/idp-api/src/internal/domain"

type SendEmailWithViewGateway interface {
	SendWithView(subject string, from string, recipients []string, views []string, layout string, data interface{}) error
}

type SendSignupEmail interface {
	Execute(value_objects.EmailAddress, string) error
}
