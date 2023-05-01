package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/lucianogarciaz/kit/vo"
	"regexp"
	"time"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type UserRepository interface {
	Update(User) error
	CreateUser(ctx context.Context, user User) error
}

type User struct {
	id        vo.ID
	createdAt time.Time
	name      string
	email     Email
	companyID vo.ID
}

func (u *User) Id() vo.ID {
	return u.id
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) CompanyID() vo.ID {
	return u.companyID
}

func New(name string, email string, companyID vo.ID) (*User, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	domainEmail, err := NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("newUser: email: %w", err)
	}

	return &User{
		id:        vo.NewID(),
		createdAt: time.Now(),
		name:      name,
		email:     domainEmail,
		companyID: companyID,
	}, nil
}

func (u *User) Hydrate(id vo.ID, createdAt time.Time, name string, email Email, companyID vo.ID) {
	u.id = id
	u.createdAt = createdAt
	u.name = name
	u.email = email
	u.companyID = companyID
}

func validateName(name string) error {
	// validate email
	return nil
}

type Email string

func NewEmail(email string) (Email, error) {
	if err := validateEmail(email); err != nil {
		return "", err
	}
	return Email(email), nil
}

var ErrEmptyEmail = errors.New("empty email")
var ErrInvalidEmail = errors.New("invalid email format")

func validateEmail(email string) error {
	if email == "" {
		return ErrEmptyEmail
	}

	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil || !matched {
		return ErrInvalidEmail
	}

	return nil
}
