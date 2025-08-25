package handler

import (
	"encoding/json"
	"net/http"

	"orders-svc/internal/kafka"
	"orders-svc/internal/models"
	"orders-svc/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service  *service.Service
	producer *kafka.Producer
}

func NewHandler(service *service.Service, producer *kafka.Producer) *Handler {
	return &Handler{
		service:  service,
		producer: producer,
	}
}

func (h *Handler) GetOrder(c *gin.Context) {
	orderUID := c.Param("id")
	if orderUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	order, source, duration, err := h.service.GetOrder(orderUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"source":   source,
		"order":    order,
		"duration": duration.String(),
	})
}

func (h *Handler) GetAllOrders(c *gin.Context) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// сериализация
	data, err := json.Marshal(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize order"})
		return
	}

	// отправка в Kafka
	if err := h.producer.Publish(c, []byte(order.OrderUID), data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish order"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Order accepted for processing",
		"orderId": order.OrderUID,
	})
}
