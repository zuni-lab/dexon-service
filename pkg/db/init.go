package db

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func Init(ctx context.Context, dbSource string, migrationURL string) {
	config, err := pgxpool.ParseConfig(dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot parse db config")
	}
	registerTypes(config)

	connPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(migrationURL, dbSource)

	DB = NewStore(connPool)
}

func registerTypes(config *pgxpool.Config) {
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		t, err := conn.LoadType(ctx, "order_type")
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(t)

		t, err = conn.LoadType(ctx, "_order_type")
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(t)

		t, err = conn.LoadType(ctx, "order_status")
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(t)

		t, err = conn.LoadType(ctx, "_order_status")
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(t)

		return err
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func Close() {
	DB.connPool.Close()
}
