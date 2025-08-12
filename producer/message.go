package producer

import (
	"encoding/json"
	"log"
	"math/rand"

	entities "github.com/muradrmagomedov/wbstestproject"
)

func wrongMessage() string {
	return "Hello world"
}

func rightMessage() string {
	order := generateOrder()
	msg, err := json.Marshal(order)
	if err != nil {
		log.Printf("error while generating message in rightMessage():%v", err.Error())
		return wrongMessage()
	}
	return string(msg)
}

func generatePayment() entities.Payment {
	payment := entities.Payment{
		RequestID:    getRandomString(6),
		Currency:     "usd",
		Provider:     getRandomString(5),
		Amount:       getRandomNumber(),
		PaymentDt:    getRandomNumber(),
		Bank:         "WB",
		DeliveryCost: getRandomNumber(),
		GoodsTotal:   getRandomNumber(),
		CustomFee:    getRandomNumber(),
	}
	return payment
}

func generateDelivery() entities.Delivery {
	delivery := entities.Delivery{
		Name:    "John Dow",
		Phone:   "+799999999999",
		Zip:     getRandomString(6),
		City:    getRandomString(10),
		Address: getRandomString(15),
		Region:  getRandomString(15),
		Email:   "123@mail.ru",
	}
	return delivery
}

func generateItem() entities.Item {
	item := entities.Item{
		ChartID:    getRandomNumber(),
		Price:      getRandomNumber(),
		RId:        getRandomString(10),
		Name:       getRandomString(10),
		Sale:       getRandomNumber(),
		Size:       getRandomString(10),
		TotalPrice: getRandomNumber(),
		NmId:       getRandomNumber(),
		Brand:      getRandomString(10),
		Status:     getRandomNumber(),
	}
	return item
}
func generateOrder() entities.Order {
	items := []entities.Item{}

	order := entities.Order{
		OrderUID:          getUID(),
		TrackNumber:       getUID(),
		Entry:             getRandomString(10),
		Delivery:          generateDelivery(),
		Payment:           generatePayment(),
		Locale:            getRandomString(3),
		InternalSignature: getRandomString(10),
		CustomerID:        getRandomString(10),
		DeliveryService:   getRandomString(10),
		ShardKey:          getRandomString(10),
		SmID:              getRandomNumber(),
		OofShard:          getRandomString(10),
	}

	for range rand.Intn(5) + 1 {
		item := generateItem()
		item.TrackNumber = order.TrackNumber
		items = append(items, item)
	}
	order.Items = items
	payment := generatePayment()
	payment.Transaction = order.OrderUID
	order.Payment = payment
	delivery := generateDelivery()
	delivery.OrderUID = order.OrderUID
	order.Delivery = delivery
	return order
}
