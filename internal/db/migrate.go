package db

import (
	"context"
	"log"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) {
	sql, err := os.ReadFile("internal/db/migrations/001_init.sql")
	if err != nil {
		log.Fatal("failed to read migration:", err)
	}

	_, err = pool.Exec(context.Background(), string(sql))
	if err != nil {
		log.Fatal("failed to run migration:", err)
	}

	log.Println("migrations applied")
}