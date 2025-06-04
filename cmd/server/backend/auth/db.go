package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateDBPool(host string, port int, name string, username string, password string, ssl string, pool_max_conns int) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s pool_max_conns=%d",
		username, password, host, port, name, ssl, pool_max_conns,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %v", err)
	}

	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}
	// Add this to anything that calls this function
	//defer dbpool.Close()
	return dbpool, nil
}
