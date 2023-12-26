package entities

import value_objects "github.com/kallabs/idp-api/src/internal/domain"

type TokenType uint8

const (
	TypeAccess TokenType = iota
	TypeRefresh
)

type TokenPayload struct {
	Type   TokenType
	UserId value_objects.ID
}

type AuthPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	// AccessExpiresIn  int
	// RefreshExpiresIn int
}
