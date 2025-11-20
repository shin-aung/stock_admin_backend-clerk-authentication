package main

import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func getenv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func initDB() error {
	    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        getenv("DATABASE_USER", "pyou_developer"),
        getenv("DATABASE_PASSWORD", "stockpassword123"),
        getenv("DATABASE_HOST", "localhost"),
        getenv("DATABASE_PORT", "5432"),
        getenv("DATABASE_NAME", "stock_dashboard"),
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
