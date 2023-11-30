package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	PoolSize int
}

func ConnectDB(cfg Config) (*pgxpool.Pool, error) {

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.PoolSize)

	return pgxpool.ConnectConfig(context.Background(), poolConfig)
}

func CloseDB(pool *pgxpool.Pool) {
	pool.Close()
}
