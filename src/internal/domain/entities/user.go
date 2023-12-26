package entities

import (
	"time"

	value_objects "github.com/kallabs/idp-api/src/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserRole ...
type UserRole uint

const (
	UserAdmin UserRole = iota
	UserReader
	UserEditor
)

// UserStatus is used to identiry user statuses
type UserStatus uint

const (
	UserInactive UserStatus = iota
	UserUnconfirmed
	UserActive
	UserDeleted
)

type MemberStatus uint

const (
	MemberPending MemberStatus = iota
	MemberActive
	MemberDeleted
)

type UserPassword struct {
	PasswordHash string `json:"-" db:"password_hash"`
}

func (up *UserPassword) SetPassword(plainPassword string) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), 14)
	up.PasswordHash = string(bytes)
}

func (up *UserPassword) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(up.PasswordHash), []byte(password))

	return err == nil
}

type User struct {
	UserPassword
	Id             *value_objects.ID          `json:"id" db:"id"`
	Email          value_objects.EmailAddress `json:"email" db:"email"`
	Username       string                     `json:"username" db:"username"`
	FirstName      string                     `json:"firstName" db:"first_name"`
	LastName       string                     `json:"lastName" db:"last_name"`
	Token          string                     `json:"-" db:"token"`
	Role           UserRole                   `json:"role" db:"role"`
	Status         UserStatus                 `json:"status" db:"status"`
	TokenExpiresAt time.Time                  `json:"-" db:"token_expires_at"`
	CreatedAt      time.Time                  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time                  `json:"updatedAt" db:"updated_at"`
}

type UserInfo struct {
	Id        *value_objects.ID          `json:"id" db:"id"`
	Email     value_objects.EmailAddress `json:"email" db:"email"`
	Username  string                     `json:"username" db:"username"`
	FirstName string                     `json:"firstName" db:"first_name"`
	LastName  string                     `json:"lastName" db:"last_name"`
	Role      UserRole                   `json:"role" db:"role"`
	Status    UserStatus                 `json:"status" db:"status"`
	CreatedAt time.Time                  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time                  `json:"updatedAt" db:"updated_at"`
}

type CreateUserSchema struct {
	UserPassword
	Email          value_objects.EmailAddress `json:"email" db:"email"`
	Username       string                     `json:"username" db:"username"`
	FirstName      string                     `json:"firstName" db:"first_name"`
	LastName       string                     `json:"lastName" db:"last_name"`
	Token          string                     `json:"token" db:"token"`
	Role           UserRole                   `json:"role" db:"role"`
	Status         UserStatus                 `json:"status" db:"status"`
	TokenExpiresAt time.Time                  `json:"-" db:"token_expires_at"`
}
