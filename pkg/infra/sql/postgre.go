package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"

	_ "github.com/lib/pq"

	"github.com/lucianogarciaz/pulley-example/pkg/domain"
)

var _ domain.UserRepository = &PostgresUserRepository{}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}
func (p PostgresUserRepository) Update(user domain.User) error {
	panic("implement me")
}

func (p PostgresUserRepository) CreateUser(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO users (id, created_at, name, email, company_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := p.db.ExecContext(ctx, query, user.Id(), user.CreatedAt(), user.Name(), user.Email(), user.CompanyID())
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return domain.ErrDuplicateEmail
			}
		}

		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}

func InitDb() (*sql.DB, error) {
	//dsn := os.Getenv("POSTGRES_DSN")
	//if dsn == "" {
	//	return nil, errors.New("missing Postgres DSN")
	//}
	dsn := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	return db, nil
}
