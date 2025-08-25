package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	KafkaURL    string
	KafkaTopic  string
	HTTPPort    string
	CacheSize   int
}

func Load() (*Config, error) {
	// Загружаем .env файл если он существует
	_ = godotenv.Load()

	cacheSizeStr := os.Getenv("CACHE_SIZE")
	cacheSize, err := strconv.Atoi(cacheSizeStr)
	if err != nil || cacheSize <= 0 {
		cacheSize = 10 // default
	}

	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://customer:customer_pass@postgres:5432/orders_db?sslmode=disable"),
		KafkaURL:    getEnv("KAFKA_URL", "kafka:9092"),
		KafkaTopic:  getEnv("KAFKA_TOPIC", "orders"),
		HTTPPort:    getEnv("HTTP_PORT", ":8080"),
		CacheSize:   cacheSize,
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
