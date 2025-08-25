package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"orders-svc/internal/models"
	"orders-svc/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader       *kafka.Reader
	dlqWriter    *kafka.Writer
	orderService *service.Service
	validate     *validator.Validate
}

func NewConsumer(brokers []string, topic string, groupID string, orderService *service.Service) *Consumer {
	// Настройка Kafka reader для получения сообщений
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		CommitInterval: 0,
		MinBytes:       1e3,
		MaxBytes:       1e6,
	})

	// Настройка DLQ writer для ошибочных сообщений
	wDLQ := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic + "-dlq",
	})

	return &Consumer{
		reader:       r,
		dlqWriter:    wDLQ,
		orderService: orderService,
		validate:     validator.New(),
	}
}

func (c *Consumer) Start() {
	for {
		ctx := context.Background()

		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error (partition=%d, offset=%d): %v\n", m.Partition, m.Offset, err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("Message received (partition=%d, offset=%d, key=%s): %s\n",
			m.Partition, m.Offset, string(m.Key), string(m.Value))

		var order models.Order

		// Декодирование JSON в структуру заказа
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Println("Invalid JSON, send to DLQ:", err)
			c.dlqWriter.WriteMessages(ctx, kafka.Message{
				Key:   m.Key,
				Value: m.Value,
			})
			continue
		}

		// Валидация структуры заказа
		if err := c.validate.Struct(order); err != nil {
			log.Println("Validation failed, send to DLQ:", err)
			c.dlqWriter.WriteMessages(ctx, kafka.Message{
				Key:   m.Key,
				Value: m.Value,
			})
			continue
		}

		// Сохранение заказа в базу данных
		if err := c.orderService.CreateOrder(&order); err != nil {
			log.Println("Failed to save order, send to DLQ:", err)
			c.dlqWriter.WriteMessages(ctx, kafka.Message{
				Key:   m.Key,
				Value: m.Value,
			})
			continue
		}

		// Подтверждение успешной обработки сообщения
		log.Printf("Order %s successfully processed\n", order.OrderUID)
		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Println("Failed to commit message:", err)
		} else {
			log.Printf("Message offset %d committed\n", m.Offset)
		}
	}
}

func (c *Consumer) Close() {
	c.reader.Close()
	c.dlqWriter.Close()
}
