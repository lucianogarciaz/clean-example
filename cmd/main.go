package main

import (
	"github.com/lucianogarciaz/kit/obs"
	"github.com/lucianogarciaz/pulley-example/pkg/app"
	"github.com/lucianogarciaz/pulley-example/pkg/infra/server"
	repo "github.com/lucianogarciaz/pulley-example/pkg/infra/sql"

	"log"
)

func main() {
	o11y := obs.NewObserver(obs.NoopMetrics{}, obs.NewBasicLogger())
	db, err := repo.InitDb()
	if err != nil {
		log.Fatal("init db:", err.Error())
	}

	ch := app.NewCreateUserCH(repo.NewPostgresUserRepository(db))
	middleware := obs.CommandHandlerObsMiddleware[app.CreateUserCommand](o11y)
	s := server.NewServer(o11y, middleware(ch))

	s.Serve()
}
