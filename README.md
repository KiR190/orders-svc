# 📦 Order Service

Микросервис для управления заказами с использованием **Go**, **PostgreSQL**, **Kafka** и кеша в памяти.  

## ✨ Возможности
- Получает сообщения с заказами из **Kafka**  
- Cохраняет их в **PostgreSQL** (транзакции)  
- Dead Letter Queue для невалидных сообщений
- Быстрый доступ к заказам через кеш (map)  
- Предоставляет **HTTP API** и простой **веб-интерфейс** для просмотра данных заказа
  

## 🛠️ Технологии
* Go 1.24.5
* Gin Web Framework
* PostgreSQL 17.6
* Apache Kafka
* Docker + Docker Compose

## 🚀 Запуск

```bash
git clone https://github.com/KiR190/orders-svc
cd orders-svc
cp .env.example .env
docker-compose up --build
```
Сервис будет доступен на: [http://localhost](http://localhost)

## 📝 Запуск скрипта-эмулятор отправки сообщений в Kafka
```bash
go run ./cmd/producer
```

## 📂 Структура проекта

```text
orders-svc/
├── cmd/
│   ├── app/main.go          # Точка входа. Основной микросервис
│   ├── producer/main.go     # Скрипт для тестовой отправки сообщений в Kafka
├── internal/
│   ├── cache/               # In-memory кеш
│   ├── database/            # Работа с PostgreSQL 
│   ├── handler/             # HTTP обработчики
│   ├── kafka/               # Consumer и Producer для Kafka
│   ├── models/              # Структуры данных
│   └── service/             # Бизнес-логика
├── scripts/
│   ├── init_db.sql          # Инициализация схемы БД, таблиц, индексов
│   └── init_users.sql       # Создание пользователей, назначение прав
└── docker-compose.yml       # Full environment
```

## 🔌 API

### Получить все заказы

```http
GET /api/orders
```

### Получить заказ по ID

```http
GET /api/order/{order_uid}
```

### Создать заказ

```http
POST /api/order
Content-Type: application/json
```

### Пример тела запроса:

```json
{
  "order_uid": "38cfe564457ec98efaa2770ftest",
  "track_number": "WBzq07pP7Z0a",
  "entry": "GOG",
  "delivery": {
    "name": "Alex Johnson",
    "phone": "+5551084406",
    "zip": "5037798",
    "city": "Night City",
    "address": "Ben Gurion 85",
    "region": "South",
    "email": "test+fOkzO@gmail.com"
  },
  "payment": {
    "transaction": "38cfe564457ec98efaa2770ftest",
    "request_id": "",
    "currency": "ILS",
    "provider": "wbpay",
    "amount": 4597,
    "payment_dt": 1744155492,
    "bank": "tbank",
    "delivery_cost": 1528,
    "goods_total": 2993,
    "custom_fee": 76
  },
  "items": [
    {
      "chrt_id": 3087370,
      "track_number": "WBzq07pP7Z0a",
      "price": 1373,
      "rid": "9b9c9887c8c2e536a669c6ed30cfbcf9test",
      "name": "Loot Box",
      "sale": 3,
      "size": "L",
      "total_price": 1331,
      "nm_id": 4721038,
      "brand": "Nintendo",
      "status": 203
    },
    {
      "chrt_id": 1402706,
      "track_number": "WBzq07pP7Z0a",
      "price": 3463,
      "rid": "a1a0dcd3cae83f863d1cd2ca3390ad98test",
      "name": "DLC Pack",
      "sale": 52,
      "size": "M",
      "total_price": 1662,
      "nm_id": 9098873,
      "brand": "Rockstar",
      "status": 201
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "cust_oFx7WY",
  "delivery_service": "ups",
  "shardkey": "33",
  "sm_id": 308,
  "date_created": "2025-04-08T23:38:12Z",
  "oof_shard": "3"
}

```

## 🛡️ Надежность и отказоустойчивость

✅ Транзакции БД - атомарное сохранение заказов

✅ Manual commit offset - контроль подтверждения Kafka

✅ DLQ - изоляция невалидных сообщений

✅ Idempotency - защита от дубликатов

## 📹 Демонстрация
Сервис в работе (видео): [ссылка на Яндекс.Диск](https://disk.yandex.ru/i/d7fxwCuozsIwWg) 