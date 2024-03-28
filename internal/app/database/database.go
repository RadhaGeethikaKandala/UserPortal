package database

import (
	"fmt"
	"log"
	"userportal/internal/app/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const postgresDatasource = "postgres://%s:%s@%s:%s/%s?sslmode=%s"
const configFileRelativePath = "./internal/app/config"

var (
	databaseConf = config.ReadConfig(configFileRelativePath).Database
	dataSource   = fmt.Sprintf(
		postgresDatasource,
		databaseConf.Username,
		databaseConf.Password,
		databaseConf.Host,
		databaseConf.Port,
		databaseConf.Dbname,
		databaseConf.Sslmode,
	)
)

func CreateDatabaseConn() *sqlx.DB {
	dbConn, err := sqlx.Open("postgres", dataSource)
	errString := "Unable to open a connection to database, error: %s"
	if err != nil {
		log.Fatalf(errString, err.Error())
	}
	if err := dbConn.Ping(); err != nil {
		log.Fatalf(errString, err.Error())
	}

	return dbConn
}

func RunDatabaseMigrations() {
	migrations_uri := "file://internal/app/database/migrations"

	m, err := migrate.New(migrations_uri, dataSource)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

}
