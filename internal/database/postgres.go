package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func Connect(databaseURL string) (*DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = databaseURL
	}

	log.Printf("Using database URL: %s", connStr)

	// Конфигурация пула соединений
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse config error: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	// Подключение с ретраями
	var pool *pgxpool.Pool
	for i := 0; i < 10; i++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			err = pool.Ping(context.Background())
			if err == nil {
				log.Println("Database connection established with pgxpool")
				return &DB{pool: pool}, nil
			}
		}
		log.Printf("Database not available, retrying... (attempt %d/10)", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("database unavailable after 10 attempts: %v", err)
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
