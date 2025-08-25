package service

import (
	"context"
	"time"

	"orders-svc/internal/cache"
	"orders-svc/internal/database"
	"orders-svc/internal/models"
)

type Service struct {
	cache     *cache.Cache
	db        *database.DB
	cacheSize int
}

func New(cache *cache.Cache, db *database.DB, cacheSize int) *Service {
	return &Service{
		cache:     cache,
		db:        db,
		cacheSize: cacheSize,
	}
}

func (s *Service) GetOrder(orderUID string) (*models.Order, string, time.Duration, error) {
	start := time.Now()
	// Сначала проверяем кеш
	if order, exists := s.cache.Get(orderUID); exists {
		elapsed := time.Since(start)
		return order, "cache", elapsed, nil
	}

	// Если нет в кеше, получаем из БД
	ctx := context.Background()
	order, err := s.db.GetOrderByUID(ctx, orderUID)
	if err != nil {
		return nil, "", 0, err
	}

	// Сохраняем в кеш для будущих запросов
	s.cache.Set(order)

	elapsed := time.Since(start)

	return order, "db", elapsed, nil
}

func (s *Service) GetAllOrders() ([]*models.Order, error) {
	ctx := context.Background()
	return s.db.GetAllOrders(ctx, s.cacheSize)
}

func (s *Service) CreateOrder(order *models.Order) error {
	ctx := context.Background()

	// Сохраняем в БД
	err := s.db.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	// Сохраняем в кеш
	s.cache.Set(order)

	return nil
}

func (s *Service) RestoreCacheFromDB(cache *cache.Cache, db *database.DB) error {
	ctx := context.Background()
	orders, err := db.GetAllOrders(ctx, s.cacheSize)
	if err != nil {
		return err
	}

	for _, order := range orders {
		cache.Set(order)
	}

	return nil
}
