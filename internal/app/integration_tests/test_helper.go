package testing

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"userportal/internal/app/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	DbInstance *sqlx.DB
	container  testcontainers.Container
}

const configFileRelativePath = "../config"

func SetupTestDatabase() *TestDatabase {

	// setup db container
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	cancel()

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
	}
}

func (tdb *TestDatabase) TearDown() {
	tdb.DbInstance.Close()
	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *sqlx.DB, error) {
	databaseConf := config.ReadConfig(configFileRelativePath).DatabaseTest

	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return nil, nil, fmt.Errorf("failed to get path")
	}

	migrationFilesPath, err := filepath.Glob(filepath.Join(filepath.Dir(path), "..", "database", "migrations", "*.sql"))

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch migration scripts from given path")
	}

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(migrationFilesPath...),
		postgres.WithDatabase(databaseConf.Dbname),
		postgres.WithUsername(databaseConf.Username),
		postgres.WithPassword(databaseConf.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return pgContainer, nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if strings.Contains(strings.ToLower(connStr), "localhost") {
		connStr = strings.Replace(connStr, "localhost", databaseConf.Host, 1)
	}

	if err != nil {
		return pgContainer, nil, err
	}
	db, err := sqlx.Open("postgres",
		connStr)

	if err != nil {
		return pgContainer, db, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return pgContainer, db, nil
}
