package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateDBPool(host string, port int, name string, username string, password string, ssl string, pool_max_conns int, logger *slog.Logger) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s pool_max_conns=%d",
		username, password, host, port, name, ssl, pool_max_conns,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Error("unable to parse config", "err", err)
		return nil, err
	}

	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		logger.Error("unable to create connection pool", "err", err)
		return nil, err
	}
	// Add this to anything that calls this function
	//defer dbpool.Close()
	logger.Debug("database pool created", "dbpool", dbpool)
	return dbpool, nil
}
