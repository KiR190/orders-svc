package main

import (
	"log"
	"orders-svc/config"
	"orders-svc/internal/cache"
	"orders-svc/internal/database"
	"orders-svc/internal/handler"
	"orders-svc/internal/kafka"
	"orders-svc/internal/service"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация кеша
	cache := cache.New()

	// Подключение к БД
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация сервиса
	orderService := service.New(cache, db, cfg.CacheSize)

	// Восстановление кеша из БД при старте
	if err := orderService.RestoreCacheFromDB(cache, db); err != nil {
		log.Printf("Failed to restore cache from DB: %v", err)
	}

	// Инициализация продюсера (для отправки сообщений в Kafaka по POST /order)
	producer := kafka.NewProducer([]string{cfg.KafkaURL}, cfg.KafkaTopic)
	defer producer.Close()

	// Запуск Kafka consumer в горутине
	go func() {
		consumer := kafka.NewConsumer([]string{cfg.KafkaURL}, cfg.KafkaTopic, "order-service-group", orderService)
		defer consumer.Close()

		consumer.Start()
	}()

	// Инициализация handler
	ginHandler := handler.NewHandler(orderService, producer)

	// Запуск HTTP сервера
	log.Printf("Starting HTTP server on %s", cfg.HTTPPort)
	if err := ginHandler.Run(cfg.HTTPPort); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
