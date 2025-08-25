package database

import (
	"context"

	"orders-svc/internal/models"
)

func (db *DB) CreateOrder(ctx context.Context, order *models.Order) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Вставляем данные заказа
	orderQuery := `INSERT INTO orders (
        order_uid, track_number, entry, locale, internal_signature,
        customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(ctx, orderQuery,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return err
	}

	// Вставляем данные доставки
	deliveryQuery := `INSERT INTO deliveries (
        order_uid, name, phone, zip, city, address, region, email
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(ctx, deliveryQuery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		return err
	}

	// Вставляем данные платежа
	paymentQuery := `INSERT INTO payments (
        order_uid, transaction, request_id, currency, provider, amount,
        payment_dt, bank, delivery_cost, goods_total, custom_fee
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(ctx, paymentQuery,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		return err
	}

	// Вставляем товары
	itemQuery := `INSERT INTO items (
        order_uid, chrt_id, track_number, price, rid, name, sale, size,
        total_price, nm_id, brand, status
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, item := range order.Items {
		_, err = tx.Exec(ctx, itemQuery,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID,
			item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
