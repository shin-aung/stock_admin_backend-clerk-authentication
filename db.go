package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func initDB() error {
	    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

    var err error
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    pool, err = pgxpool.New(ctx, dsn)
    if err != nil {
        return err
    }
    return nil
}
