package app_test

import (
	"context"
	"errors"
	"github.com/lucianogarciaz/kit/vo"
	"testing"

	"github.com/lucianogarciaz/pulley-example/pkg/app"
	"github.com/lucianogarciaz/pulley-example/pkg/domain"
	"github.com/stretchr/testify/require"
)

func TestCreateUserHandle(t *testing.T) {
	require := require.New(t)

	t.Run(`Given a valid CreateUserCommand, 
		when the command is handled, 
		then it should create the user`, func(t *testing.T) {
		name := "John Doe"
		email := "john@example.com"
		companyID := "d3cf153e-974f-4f93-b7e9-318bde3c1a9c"

		userRepoMock := &UserRepositoryMock{
			CreateUserFunc: func(ctx context.Context, user domain.User) error {
				return nil
			},
		}

		createUserCH := app.NewCreateUserCH(userRepoMock)
		cmd := app.CreateUserCommand{Name: name, Email: email, CompanyID: companyID}

		events, err := createUserCH.Handle(context.Background(), cmd)

		require.NoError(err)
		require.Nil(events)
	})
	t.Run(`Given an invalid companyID in CreateUserCommand, 
	when the command is handled, 
	then it should return an error`, func(t *testing.T) {
		name := "John Doe"
		email := "john@example.com"
		invalidCompanyID := "invalid-company-id"

		userRepoMock := &UserRepositoryMock{}

		createUserCH := app.NewCreateUserCH(userRepoMock)
		cmd := app.CreateUserCommand{Name: name, Email: email, CompanyID: invalidCompanyID}

		events, err := createUserCH.Handle(context.Background(), cmd)

		require.ErrorIs(err, vo.ErrInvalidID)
		require.Nil(events)
	})

	t.Run(`Given an invalid email in CreateUserCommand, 
	when the command is handled, 
	then it should return an error`, func(t *testing.T) {
		name := "John Doe"
		invalidEmail := "john@example"
		companyID := "d3cf153e-974f-4f93-b7e9-318bde3c1a9c"

		userRepoMock := &UserRepositoryMock{}

		createUserCH := app.NewCreateUserCH(userRepoMock)
		cmd := app.CreateUserCommand{Name: name, Email: invalidEmail, CompanyID: companyID}

		events, err := createUserCH.Handle(context.Background(), cmd)

		require.ErrorIs(err, domain.ErrInvalidEmail)
		require.Nil(events)
	})

	t.Run(`Given a valid CreateUserCommand but a repository that returns error, 
		when the command is handled, 
		then it should return an error`, func(t *testing.T) {
		name := "John Doe"
		email := "john@example.com"
		companyID := "d3cf153e-974f-4f93-b7e9-318bde3c1a9c"
		expectedError := errors.New("addUser")

		userRepoMock := &UserRepositoryMock{
			CreateUserFunc: func(ctx context.Context, user domain.User) error {
				return expectedError
			},
		}

		createUserCH := app.NewCreateUserCH(userRepoMock)
		cmd := app.CreateUserCommand{Name: name, Email: email, CompanyID: companyID}

		events, err := createUserCH.Handle(context.Background(), cmd)

		require.ErrorIs(err, expectedError)
		require.Nil(events)
	})
}
