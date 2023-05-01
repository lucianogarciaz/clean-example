package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lucianogarciaz/kit/vo"
	"github.com/lucianogarciaz/pulley-example/pkg/app"
	"github.com/lucianogarciaz/pulley-example/pkg/domain"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/lucianogarciaz/kit/cqs"
	"github.com/lucianogarciaz/kit/obs"
)

const (
	timeout     = 15
	idleTimeout = 60
)

type Server struct {
	obs obs.Observer
	mux.Router
	createUser cqs.CommandHandler[app.CreateUserCommand]
}

func NewServer(obs obs.Observer,
	ch cqs.CommandHandler[app.CreateUserCommand],
) *Server {
	return &Server{obs: obs, createUser: ch}
}

func (s *Server) Serve() {
	router := mux.NewRouter()

	router.HandleFunc("/user", s.CreateUser()).Methods(http.MethodPost)

	router.Use(s.logMiddleware)
	router.Use(s.jsonContentTypeMiddleware)

	server := &http.Server{
		Addr:         port(),
		Handler:      router,
		ReadTimeout:  timeout * time.Second,
		WriteTimeout: timeout * time.Second,
		IdleTimeout:  idleTimeout * time.Second,
	}

	_ = s.obs.Log(obs.LevelInfo, fmt.Sprintf("http server Listening on port %s", port()))

	log.Fatal(server.ListenAndServe())
}

func (s *Server) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateUserRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cmd := app.CreateUserCommand{
			Name:      request.Name,
			Email:     request.Email,
			CompanyID: request.CompanyID,
		}
		var statusCode int
		_, err = s.createUser.Handle(r.Context(), cmd)
		if err != nil {
			statusCode = http.StatusInternalServerError
			if errors.Is(err, domain.ErrInvalidEmail) ||
				errors.Is(err, domain.ErrEmptyEmail) ||
				errors.Is(err, domain.ErrDuplicateEmail) ||
				errors.Is(err, vo.ErrInvalidID) {
				statusCode = http.StatusBadRequest
			}
			w.WriteHeader(statusCode)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) jsonContentTypeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func (s *Server) logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			elapsed := time.Since(begin)
			_ = s.obs.Log(obs.LevelInfo, fmt.Sprintf("request latency, %f", elapsed.Seconds()))
		}(time.Now())

		_ = s.obs.Log(
			obs.LevelInfo,
			fmt.Sprintf("new request with method %s to url: %s", r.Method, r.URL.String()),
		)
		h.ServeHTTP(w, r)
	})
}

const defaultAddress = "3001"

func port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultAddress
	}

	return fmt.Sprintf(":%s", port)
}

type CreateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	CompanyID string `json:"company_id"`
}
