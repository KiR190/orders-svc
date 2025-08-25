package cache

import (
	"orders-svc/internal/models"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]*models.Order
}

func New() *Cache {
	return &Cache{
		items: make(map[string]*models.Order),
	}
}

func (c *Cache) Get(orderUID string) (*models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.items[orderUID]
	return order, exists
}

func (c *Cache) Set(order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[order.OrderUID] = order
}

func (c *Cache) GetAll() map[string]*models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.items
}
