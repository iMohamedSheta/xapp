package load

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/database"
	"github.com/imohamedsheta/xapp/app/database/seeders"
	"github.com/imohamedsheta/xapp/app/domain/hooks"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xioc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func InitDatabase(c *xioc.Container) {
	err := xioc.Singleton(c, func(c *xioc.Container) (*xqb.DBM, error) {
		cfg := x.Config()
		defaultConnection := cfg.GetString("database.default", "default")
		dialect := cfg.GetString("database.connections."+defaultConnection+".dialect", "")
		databaseCfg := cfg.GetMap("database.connections."+defaultConnection, nil)
		dsn := utils.BuildPostgresDSN(databaseCfg)

		config, err := pgx.ParseConfig(dsn)
		if err != nil {
			return nil, err
		}

		db := stdlib.OpenDB(*config)
		configureDatabaseConnectionPool(db)
		if err := db.PingContext(context.Background()); err != nil {
			return nil, fmt.Errorf("database connection failed: %w", err)
		}

		// Add connection first
		_ = xqb.AddConnection(&xqb.Connection{
			Name:    "default",
			DB:      db,
			Dialect: "postgres",
		})

		hooks.InitQueryHooks()

		utils.Dump("Database migrations completed successfully")

		// Run migrations
		if err := runMigrations(db, dialect); err != nil {
			x.Logger().Error("Failed to run migrations: " + err.Error())
			utils.PrintErr("Failed to run migrations: " + err.Error())
		}

		// Run seeds
		// It will fail if already there the seeds
		if err := runSeeds(); err != nil {
			x.Logger().Error("Failed to run seeds: " + err.Error())
			utils.PrintErr("Failed to run seeds: " + err.Error())
		}

		return xqb.DBManager(), nil
	})

	if err != nil {
		x.Logger().Error("Failed to load xqb module in the ioc container: " + err.Error())
	}

}

func configureDatabaseConnectionPool(db *sql.DB) {
	maxIdleConns := x.Config().GetInt("database.max_idle_conns", 10)
	maxOpenConns := x.Config().GetInt("database.max_open_conns", 100)
	connMaxLifetime := x.Config().GetDuration("database.conn_max_lifetime", 30*time.Minute)

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)
}

func runMigrations(db *sql.DB, dialect string) error {
	goose.SetBaseFS(database.MigrationsFS)

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	x.Logger().Info("Database migrations completed successfully")
	return nil
}

func runSeeds() error {
	return seeders.SeedDatabase()
}
