package server_test

import (
	"bytes"
	"context"
	"github.com/lucianogarciaz/kit/vo"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lucianogarciaz/kit/cqs"
	"github.com/lucianogarciaz/pulley-example/pkg/app"
	"github.com/lucianogarciaz/pulley-example/pkg/domain"
	"github.com/lucianogarciaz/pulley-example/pkg/infra/server"
)

func TestCreateUser(t *testing.T) {
	require := require.New(t)

	t.Run("Given a valid CreateUser request, when the endpoint is called, then it should return a 200 status", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{
			"name": "John Doe",
			"email": "john.doe@example.com",
			"company_id": "3b6127f7-6d7f-4121-a634-35f27b7f3afc"
		}`))

		chMock := &CommandHandlerMock[app.CreateUserCommand]{
			HandleFunc: func(context.Context, app.CreateUserCommand) ([]cqs.Event, error) {
				return nil, nil
			},
		}

		s := server.NewServer(nil, chMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/user", body)

		s.CreateUser().ServeHTTP(w, r)

		require.Equal(http.StatusOK, w.Code)
		require.Len(chMock.HandleCalls(), 1)

	})

	t.Run("Given a CreateUser request with an empty email, when the endpoint is called, then it should return a 400 status", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{
			"name": "John Doe",
			"email": "",
			"company_id": "3b6127f7-6d7f-4121-a634-35f27b7f3afc"
		}`))

		chMock := &CommandHandlerMock[app.CreateUserCommand]{
			HandleFunc: func(context.Context, app.CreateUserCommand) ([]cqs.Event, error) {
				return nil, domain.ErrEmptyEmail
			},
		}

		s := server.NewServer(nil, chMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/user", body)
		s.CreateUser().ServeHTTP(w, r)

		require.Equal(http.StatusBadRequest, w.Code)
		require.Len(chMock.HandleCalls(), 1)
	})

	t.Run("Given a CreateUser request with an invalid email, when the endpoint is called, then it should return a 400 status", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{
			"name": "John Doe",
			"email": "invalid-email",
			"company_id": "3b6127f7-6d7f-4121-a634-35f27b7f3afc"
		}`))

		chMock := &CommandHandlerMock[app.CreateUserCommand]{
			HandleFunc: func(context.Context, app.CreateUserCommand) ([]cqs.Event, error) {
				return nil, domain.ErrInvalidEmail
			},
		}

		s := server.NewServer(nil, chMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/user", body)
		s.CreateUser().ServeHTTP(w, r)

		require.Equal(http.StatusBadRequest, w.Code)
		require.Len(chMock.HandleCalls(), 1)
	})

	t.Run("Given a CreateUser request with an invalid company ID, when the endpoint is called, then it should return a 400 status", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{
			"name": "John Doe",
			"email": "john.doe@example.com",
			"company_id": "invalid"
		}`))

		chMock := &CommandHandlerMock[app.CreateUserCommand]{
			HandleFunc: func(context.Context, app.CreateUserCommand) ([]cqs.Event, error) {
				return nil, vo.ErrInvalidID
			},
		}

		s := server.NewServer(nil, chMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/user", body)
		s.CreateUser().ServeHTTP(w, r)

		require.Equal(http.StatusBadRequest, w.Code)
		require.Len(chMock.HandleCalls(), 1)
	})

	t.Run("Given a CreateUser request with a JSON syntax error, when the endpoint is called, then it should return a 400 status", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{
			"name": "John Doe",
			"email": "john.doe@example.com",
			"company_id": "3b6127f7-6d7f-4121-a634-35f27b7f3afc"
		`))

		chMock := &CommandHandlerMock[app.CreateUserCommand]{
			HandleFunc: func(context.Context, app.CreateUserCommand) ([]cqs.Event, error) {
				return nil, nil
			},
		}

		s := server.NewServer(nil, chMock)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/user", body)
		s.CreateUser().ServeHTTP(w, r)

		require.Equal(http.StatusBadRequest, w.Code)
		require.Len(chMock.HandleCalls(), 0)
	})
}
