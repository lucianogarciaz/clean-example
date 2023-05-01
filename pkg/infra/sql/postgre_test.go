package sql_test

import (
	"context"
	"github.com/lucianogarciaz/kit/vo"
	"github.com/lucianogarciaz/pulley-example/pkg/domain"
	"github.com/lucianogarciaz/pulley-example/pkg/infra/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPostgresUserRepository(t *testing.T) {
	db, err := sql.InitDb()
	require.NoError(t, err)
	repo := sql.NewPostgresUserRepository(db)

	ctx := context.Background()

	t.Run("Create User", func(t *testing.T) {
		var user domain.User
		user.Hydrate(
			vo.NewID(),
			time.Now(),
			"John Doe",
			"john.doe@example.com",
			vo.NewID(),
		)

		// Test adding a user
		err := repo.CreateUser(ctx, user)
		require.NoError(t, err)
	})

	t.Run("error on duplicate email", func(t *testing.T) {
		var user domain.User
		user.Hydrate(
			vo.NewID(),
			time.Now(),
			"Luciano GZ",
			"luciano@pulley.com",
			vo.NewID(),
		)
		err = repo.CreateUser(ctx, user)
		require.NoError(t, err)
		var userWithDuplicateEmail domain.User
		userWithDuplicateEmail.Hydrate(
			vo.NewID(),
			time.Now(),
			"Luciano GZ2",
			"luciano@pulley.com",
			vo.NewID(),
		)

		err = repo.CreateUser(ctx, userWithDuplicateEmail)
		require.ErrorIs(t, err, domain.ErrDuplicateEmail)
	})
}
