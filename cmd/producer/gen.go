package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
	"strconv"
	"time"
)

/*func generator() {
	mrand.Seed(time.Now().UnixNano())

	order := generateOrder()
	out, err := json.MarshalIndent(order, "", "   ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}*/

func generateOrder() Order {
	// Общие справочники для рандомизации
	names := []string{"John Doe", "Alex Johnson", "Maria Petrova", "Chen Li", "Sara Cohen"}
	cities := []string{"Night City", "Vice City", "Silent Hill", "Gothem", "Midgar"}
	regions := []string{"Kraiot", "Center", "North", "South", "Coastal"}
	streetNames := []string{"Herzl", "Ben Gurion", "Ploshad Mira", "Rothschild", "King George"}
	providers := []string{"wbpay", "paypal", "stripe", "sbp"}
	banks := []string{"alpha", "sber", "tbank", "vtb"}
	currencies := []string{"USD", "EUR", "ILS", "RUB"}
	entryCodes := []string{"WBIL", "STEAM", "EPIC", "GOG"}
	services := []string{"meest", "dhl", "fedex", "ups"}

	track := "WB" + randAlphaNum(10)
	transaction := randHex(12) + "test"
	customerID := "cust_" + randAlphaNum(6)

	// Дата создания: случайная в промежутке последних 3 лет
	/*now := time.Now().UTC()
	past := now.AddDate(-3, 0, 0).Unix()
	dateCreated := time.Unix(randRangeInt64(past, now.Unix()), 0).UTC()*/

	// текущая дата
	dateCreated := time.Now().UTC().Truncate(time.Second)

	// Генерируем товары
	itemsCount := randRangeInt(1, 5)
	items := make([]Item, 0, itemsCount)
	brands := []string{"Bethesda", "Ubisoft", "Nintendo", "Rockstar", "CD Projekt"}
	itemNames := []string{"Collector's Edition", "DLC Pack", "In-Game Skin", "Loot Box", "Season Pass"}
	sizes := []string{"0", "S", "M", "L"}
	statuses := []int{200, 201, 202, 203, 210, 301}

	goodsTotal := 0
	for i := 0; i < itemsCount; i++ {
		price := randRangeInt(100, 5000)
		sale := randRangeInt(0, 70)
		total := price * (100 - sale) / 100

		it := Item{
			ChrtID:      randRangeInt(1_000_000, 9_999_999),
			TrackNumber: track,
			Price:       price,
			Rid:         randHex(16) + "test",
			Name:        itemNames[mrand.Intn(len(itemNames))],
			Sale:        sale,
			Size:        sizes[mrand.Intn(len(sizes))],
			TotalPrice:  total,
			NmID:        randRangeInt(1_000_000, 9_999_999),
			Brand:       brands[mrand.Intn(len(brands))],
			Status:      statuses[mrand.Intn(len(statuses))],
		}
		items = append(items, it)
		goodsTotal += total
	}

	deliveryCost := randRangeInt(0, 3000)
	customFee := randRangeInt(0, 200)
	amount := goodsTotal + deliveryCost + customFee

	order := Order{
		OrderUID:    transaction,
		TrackNumber: track,
		Entry:       entryCodes[mrand.Intn(len(entryCodes))],
		Delivery: Delivery{
			Name:    names[mrand.Intn(len(names))],
			Phone:   "+" + strconv.Itoa(randRangeInt(1, 9)) + strconv.Itoa(randRangeInt(100000000, 999999999)),
			Zip:     strconv.Itoa(randRangeInt(1000000, 9999999)),
			City:    cities[mrand.Intn(len(cities))],
			Address: fmt.Sprintf("%s %d", streetNames[mrand.Intn(len(streetNames))], randRangeInt(1, 200)),
			Region:  regions[mrand.Intn(len(regions))],
			Email:   "test+" + randAlphaNum(5) + "@gmail.com",
		},
		Payment: Payment{
			Transaction:  transaction,
			RequestID:    "",
			Currency:     currencies[mrand.Intn(len(currencies))],
			Provider:     providers[mrand.Intn(len(providers))],
			Amount:       amount,
			PaymentDt:    dateCreated.Unix(),
			Bank:         banks[mrand.Intn(len(banks))],
			DeliveryCost: deliveryCost,
			GoodsTotal:   goodsTotal,
			CustomFee:    customFee,
		},
		Items:             items,
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        customerID,
		DeliveryService:   services[mrand.Intn(len(services))],
		Shardkey:          strconv.Itoa(randRangeInt(0, 99)),
		SmID:              randRangeInt(1, 500),
		DateCreated:       dateCreated,
		OofShard:          strconv.Itoa(randRangeInt(1, 3)),
	}

	return order
}

func randRangeInt(min, max int) int {
	if max <= min {
		return min
	}
	return min + mrand.Intn(max-min+1)
}

func randRangeInt64(min, max int64) int64 {
	if max <= min {
		return min
	}
	return min + mrand.Int63n(max-min+1)
}

func randHex(nBytes int) string {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		// fallback на math/rand в крайне редком случае
		for i := range b {
			b[i] = byte(mrand.Intn(256))
		}
	}
	return hex.EncodeToString(b)
}

func randAlphaNum(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = letters[mrand.Intn(len(letters))]
	}
	return string(out)
}
