package services

import "github.com/kallabs/sso-api/src/internal/app/valueobject"

type EmailService interface {
	SendSignup(valueobject.EmailAddress, string) error
}
