package config

import (
	"context"
	"fmt"
	"os"

	"github.com/KylerJacobson/Go-Blog-API/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDBConn() *pgxpool.Pool {

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, db)

	logger.Sugar.Infof("Trying to connect to database %s", connStr)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return pool
}
