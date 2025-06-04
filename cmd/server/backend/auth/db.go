package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectDb() *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("GOKB_DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'database connected successfully'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

	// Add this to anything that calls this function
	//defer dbpool.Close()
	return dbpool
}
