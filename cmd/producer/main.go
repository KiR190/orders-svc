package main

import (
	"context"
	"encoding/json"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9094"},
		Topic:   "orders",
	})
	defer writer.Close()

	order := generateOrder()

	msg, err := json.MarshalIndent(order, "", "   ")
	if err != nil {
		log.Fatal("Ошибка при маршалинге:", err)
	}

	err = writer.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
	if err != nil {
		log.Printf("Ошибка при отправке: %v\n", err)
	} else {
		log.Printf("Сообщение отправлено: %s\n", string(msg))
	}
}
