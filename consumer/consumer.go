package consumer

import (
	"context"
	"encoding/json"
	"log"

	entities "github.com/muradrmagomedov/wbstestproject"
	"github.com/muradrmagomedov/wbstestproject/storage/postgresqldb"
	"github.com/segmentio/kafka-go"
)

var (
	topic     = "my-topic"
	partition = 0
)

type Consumer struct{}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c Consumer) ReadMessage(db postgresqldb.PostgresqlDB) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topic,
		Partition: partition,
	})
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		order := entities.Order{}
		err = json.Unmarshal(m.Value, &order)
		if err != nil {
			log.Printf("error while parsing order message from kafka:%v", err.Error())
			continue
		}
		db.AddOrder(order, context.Background())
	}
	if err := reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
