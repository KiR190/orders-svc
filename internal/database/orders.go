package database

import (
	"context"
	"fmt"

	"orders-svc/internal/models"
)

func (db *DB) GetAllOrders(ctx context.Context, limit int) ([]*models.Order, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `SELECT 
		order_uid, track_number, entry, locale, internal_signature, 
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard 
		FROM orders
		ORDER BY date_created DESC
		LIMIT $1`

	rows, err := db.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var orders []*models.Order

	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.OrderUID,
			&o.TrackNumber,
			&o.Entry,
			&o.Locale,
			&o.InternalSignature,
			&o.CustomerID,
			&o.DeliveryService,
			&o.Shardkey,
			&o.SmID,
			&o.DateCreated,
			&o.OofShard,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		// Загружаем связанные данные
		if err := db.loadOrderRelations(ctx, &o); err != nil {
			return nil, err
		}

		orders = append(orders, &o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return orders, nil
}

func (db *DB) GetOrderByUID(ctx context.Context, orderUID string) (*models.Order, error) {
	query := `SELECT 
		order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1`

	var order models.Order
	err := db.pool.QueryRow(ctx, query, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		return nil, fmt.Errorf("order query error: %v", err)
	}

	// Загружаем связанные данные
	if err := db.loadOrderRelations(ctx, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (db *DB) loadOrderRelations(ctx context.Context, order *models.Order) error {
	// Delivery
	var delivery models.Delivery
	deliveryQuery := `SELECT name, phone, zip, city, address, region, email FROM deliveries WHERE order_uid = $1`
	err := db.pool.QueryRow(ctx, deliveryQuery, order.OrderUID).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("delivery query error: %v", err)
	}
	order.Delivery = delivery

	// Payment
	var payment models.Payment
	paymentQuery := `SELECT 
		transaction, request_id, currency, provider, amount,
		payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments WHERE order_uid = $1`
	err = db.pool.QueryRow(ctx, paymentQuery, order.OrderUID).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("payment query error: %v", err)
	}
	order.Payment = payment

	// Items
	itemsQuery := `SELECT 
		chrt_id, track_number, price, rid, name, sale, size,
		total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1`
	rows, err := db.pool.Query(ctx, itemsQuery, order.OrderUID)
	if err != nil {
		return fmt.Errorf("items query error: %v", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return fmt.Errorf("item scan error: %v", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("items rows error: %v", err)
	}

	order.Items = items
	return nil
}
