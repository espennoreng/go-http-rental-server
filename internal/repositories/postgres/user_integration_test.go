package postgres_test

import (
	"context"
	"testing"

	"github.com/espennoreng/go-http-rental-server/internal/repositories"
	repoPostgres "github.com/espennoreng/go-http-rental-server/internal/repositories/postgres"
	"github.com/stretchr/testify/require"
)

func TestPostgresUserRepository(t *testing.T) {
	// Truncate all tables that might be affected.
	// Use CASCADE to handle foreign key relationships correctly.
	_, err := dbpool.Exec(context.Background(), "TRUNCATE users RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	repo := repoPostgres.NewUserRepository(dbpool)

	repositories.UserRepositoryTestSuite(t, repo)
}
