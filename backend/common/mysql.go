package common

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"liveroom-battle/config"
)

func InitMySQL(cfg config.MySQLConfig) *sql.DB {
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		slog.Error("failed to open mysql", "err", err)
		os.Exit(1)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("failed to connect to mysql", "err", err)
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("mysql connected: %s", cfg.DSN))
	return db
}
