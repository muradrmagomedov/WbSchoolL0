package producer

import (
	"encoding/json"
	"log"
	"math/rand"

	faker "github.com/brianvoe/gofakeit/v7"
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
		RequestID:    faker.LetterN(6),
		Currency:     faker.Currency().Short,
		Provider:     faker.LetterN(5),
		Amount:       faker.Number(0, 100000),
		PaymentDt:    faker.Number(0, 100000),
		Bank:         faker.BankName(),
		DeliveryCost: faker.Number(0, 100000),
		GoodsTotal:   faker.Number(0, 100000),
		CustomFee:    faker.Number(0, 100000),
	}
	return payment
}

func generateDelivery() entities.Delivery {
	delivery := entities.Delivery{
		Name:    faker.Name(),
		Phone:   faker.Phone(),
		Zip:     faker.Zip(),
		City:    faker.City(),
		Address: faker.Address().Address,
		Region:  faker.Address().Country,
		Email:   faker.Email(),
	}
	return delivery
}

func generateItem() entities.Item {
	item := entities.Item{
		ChartID:    faker.Number(0, 100000),
		Price:      faker.Number(0, 100000),
		RId:        faker.LetterN(10),
		Name:       faker.ProductName(),
		Sale:       faker.Number(0, 1000),
		Size:       faker.LetterN(10),
		TotalPrice: faker.Number(0, 1000),
		NmId:       faker.Number(0, 100000),
		Brand:      faker.LetterN(10),
		Status:     faker.Number(0, 10),
	}
	return item
}
func generateOrder() entities.Order {
	items := []entities.Item{}

	order := entities.Order{
		OrderUID:          faker.LetterN(15),
		TrackNumber:       faker.LetterN(15),
		Entry:             faker.LetterN(10),
		Delivery:          generateDelivery(),
		Payment:           generatePayment(),
		Locale:            faker.LetterN(3),
		InternalSignature: faker.LetterN(10),
		CustomerID:        faker.LetterN(10),
		DeliveryService:   faker.LetterN(10),
		ShardKey:          faker.LetterN(10),
		SmID:              faker.Number(0, 100000),
		OofShard:          faker.LetterN(10),
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
