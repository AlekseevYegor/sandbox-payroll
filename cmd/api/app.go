package main

import (
	"database/sql"
	"fmt"
	"se-challenge-payroll/internal/api"
	tserv "se-challenge-payroll/internal/domain/timetracking"
	"se-challenge-payroll/internal/repository"
	"se-challenge-payroll/internal/service"
	"se-challenge-payroll/pkg/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func newService() *service.Service {

	db := connectDB()

	ts := tserv.New(repository.NewTimeTrackerRepo(db))

	pr := repository.NewPayrollRepo(db)

	svc := service.NewService(ts, pr)

	return &svc

}

func connectDB() *sqlx.DB {
	var (
		host     = viper.GetString("DB_HOST")
		port     = viper.GetString("DB_PORT")
		dbName   = viper.GetString("DB_NAME")
		user     = viper.GetString("DB_USER")
		password = viper.GetString("DB_PASSWORD")
		schema   = viper.GetString("DB_SCHEMA")
		ssl      = viper.GetBool("SSL_MODE")
	)

	address := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?search_path=%s&sslmode=", user, password, host, port, dbName, schema)
	if !ssl {
		address += "disable"
	} else {
		address += "require"
	}

	db, err := sql.Open(api.DBDriverName, address)
	if err != nil {
		log.ZL.Fatal().Err(err).Msg("failed to connect to postgres DB")
	}
	if err := db.Ping(); err != nil {
		log.ZL.Fatal().Err(err).Msg("failed to ping to postgres DB")
	}

	log.ZL.Info().Msg("checking DB migrations")

	if err := applySchemaMigrationWithDatabaseInstance(api.DBDriverName, db); err != nil {
		log.ZL.Fatal().Err(err).Msg("failed to migrate postgres DB schema")
	}

	log.ZL.Info().Msg("DB connected")

	return sqlx.NewDb(db, api.DBDriverName)
}

// applySchemaMigrationWithDatabaseInstance creates and migrates db schema versions; based on db driver instance
// Doesn't close db connection though this is fully caller's responsibility
func applySchemaMigrationWithDatabaseInstance(databaseName string, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", databaseName, driver)
	if err != nil {
		return err
	}

	fmt.Println("Before Up")

	err = m.Up()
	if err != migrate.ErrNoChange {
		return err
	}

	_, _, err = m.Version()
	return err

}
