package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"

	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/marco-almeida/challenge-02/internal/config"
	"github.com/marco-almeida/challenge-02/internal/postgresql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	// get env vars
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	setupLogging(config)

	// setup graceful shutdown signals
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	// init db
	connPool, err := postgresql.NewPostgreSQL(ctx, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	// run migrations
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresHost, config.PostgresPort, config.PostgresDatabase)

	err = runDBMigration(config.MigrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot run db migration")
	}

	log.Info().Msg("db migrated successfully")

	// running in waitgroup coroutine in order to wait for graceful shutdown
	waitGroup, ctx := errgroup.WithContext(ctx)

	runHTPPServer(ctx, waitGroup, config, connPool)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runHTPPServer(ctx context.Context, waitGroup *errgroup.Group, config config.Config, connPool *pgxpool.Pool) {
	
}

func runDBMigration(migrationURL string, dbSource string) error {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		return fmt.Errorf("cannot create new migrate instance: %w", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrate up: %w", err)
	}

	return nil
}

func setupLogging(config config.Config) {
	// log to file ./logs/challenge-02/main.log and terminal
	logFolder := filepath.Join("logs", "challenge-02")
	err := os.MkdirAll(logFolder, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create log folder")
	}

	logFile := filepath.Join(logFolder, "main.log")

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot open log file")
	}

	// set up json or human readable logging
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: io.MultiWriter(os.Stdout, f)})
	} else {
		log.Logger = log.Output(io.MultiWriter(os.Stdout, f))
	}
}