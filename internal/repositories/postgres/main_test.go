package postgres_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	repoPostgres "github.com/espennoreng/go-http-rental-server/internal/repositories/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// dbpool is a package-level variable that will be accessible
// by all test files in the postgres_test package.
var dbpool *pgxpool.Pool

// TestMain is the single entry point for all tests in this package.
// It sets up the database container and connection pool.
func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start PostgreSQL container: %v", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate PostgreSQL container: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	migrator, err := migrate.New("file://../../../migrations", connStr)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	dbpool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
	defer dbpool.Close()

	// Run all tests in the package
	code := m.Run()

	os.Exit(code)
}

type TestHelper struct {
	dbpool      *pgxpool.Pool
	orgRepo     *repoPostgres.OrganizationRepository
	userRepo    *repoPostgres.UserRepository
	orgUserRepo *repoPostgres.OrganizationUserRepository
}

func SetupTestHelper(t *testing.T) *TestHelper {
	ctx := context.Background()
	_, err := dbpool.Exec(ctx, "TRUNCATE organizations RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	return &TestHelper{
		dbpool:      dbpool,
		orgRepo:     repoPostgres.NewOrganizationRepository(dbpool),
		userRepo:    repoPostgres.NewUserRepository(dbpool),
		orgUserRepo: repoPostgres.NewOrganizationUserRepository(dbpool),
	}
}

func (th *TestHelper) ResetDB(t *testing.T) {
	ctx := context.Background()
	_, err := th.dbpool.Exec(ctx, "TRUNCATE organizations RESTART IDENTITY CASCADE")
	require.NoError(t, err)
	_, err = th.dbpool.Exec(ctx, "TRUNCATE users RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}
