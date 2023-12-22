package services

import "github.com/kallabs/idp-api/src/internal/app/valueobject"

type EmailService interface {
	SendSignup(valueobject.EmailAddress, string) error
}
