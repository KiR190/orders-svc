package main

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required"`        // обязательный
	TrackNumber       string    `json:"track_number" validate:"required"`     // обязательный
	Entry             string    `json:"entry"`                                // опциональный
	Delivery          Delivery  `json:"delivery" validate:"required"`         // обязательный
	Payment           Payment   `json:"payment" validate:"required"`          // обязательный
	Items             []Item    `json:"items" validate:"required,min=1"`      // должен быть хотя бы 1 товар
	Locale            string    `json:"locale"`                               // опциональный
	InternalSignature string    `json:"internal_signature"`                   // опциональный
	CustomerID        string    `json:"customer_id" validate:"required"`      // обязательный
	DeliveryService   string    `json:"delivery_service" validate:"required"` // обязательный
	Shardkey          string    `json:"shardkey"`                             // опциональный
	SmID              int       `json:"sm_id" validate:"gte=0"`               // >= 0
	DateCreated       time.Time `json:"date_created" validate:"required"`     // обязательный
	OofShard          string    `json:"oof_shard"`                            // опциональный
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`        // обязательный
	Phone   string `json:"phone" validate:"required,e164"`  // телефоный номер
	Zip     string `json:"zip" validate:"required"`         // обязательный
	City    string `json:"city" validate:"required"`        // обязательный
	Address string `json:"address" validate:"required"`     // обязательный
	Region  string `json:"region"`                          // опциональный
	Email   string `json:"email" validate:"required,email"` // email
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`    // обязательный
	RequestID    string `json:"request_id"`                         // опциональный
	Currency     string `json:"currency" validate:"required,len=3"` // 3 символа
	Provider     string `json:"provider" validate:"required"`       // обязательный
	Amount       int    `json:"amount" validate:"gte=0"`            // >= 0
	PaymentDt    int64  `json:"payment_dt" validate:"gte=0"`        // >= 0
	Bank         string `json:"bank"`                               // опциональный
	DeliveryCost int    `json:"delivery_cost" validate:"gte=0"`     // >= 0
	GoodsTotal   int    `json:"goods_total" validate:"gte=0"`       // >= 0
	CustomFee    int    `json:"custom_fee" validate:"gte=0"`        // >= 0
}

type Item struct {
	ChrtID      int    `json:"chrt_id" validate:"required"`      // обязательный
	TrackNumber string `json:"track_number" validate:"required"` // обязательный
	Price       int    `json:"price" validate:"gte=0"`           // >= 0
	Rid         string `json:"rid" validate:"required"`          // обязательный
	Name        string `json:"name" validate:"required"`         // обязательный
	Sale        int    `json:"sale" validate:"gte=0,lte=100"`    // скидка от 0 до 100%
	Size        string `json:"size"`                             // опциональный
	TotalPrice  int    `json:"total_price" validate:"gte=0"`     // >= 0
	NmID        int    `json:"nm_id" validate:"required"`        // обязательный
	Brand       string `json:"brand"`                            // опциональный
	Status      int    `json:"status" validate:"gte=0"`          // >= 0
}
