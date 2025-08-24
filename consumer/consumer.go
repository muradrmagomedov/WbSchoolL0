package consumer

import (
	"context"
	"encoding/json"
	"log"

	entities "github.com/muradrmagomedov/wbstestproject"
	"github.com/muradrmagomedov/wbstestproject/storage/postgresql"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type KafkaConsumer struct {
	Conn *kafka.Reader
}

func NewConsumer() *KafkaConsumer {
	return &KafkaConsumer{}
}

type ConsumerConfig struct {
	Addr      string
	Topic     string
	Partition int
	GroupID   string
}

func createConfig() ConsumerConfig {
	return ConsumerConfig{
		Addr:      viper.GetString("kafka.addr"),
		Topic:     viper.GetString("kafka.topic"),
		Partition: viper.GetInt("kafka.partition"),
		GroupID:   viper.GetString("kafka.group_ud"),
	}
}

func (c *KafkaConsumer) ReadMessage(db postgresql.DB) {
	config := createConfig()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{config.Addr},
		Topic:     config.Topic,
		Partition: config.Partition,
		GroupID:   config.GroupID,
	})
	c.Conn = reader
	ctx := context.Background()
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			break
		}
		log.Printf("offset of message:%v", m.Offset)
		err = reader.CommitMessages(ctx, m)
		if err != nil {
			log.Printf("couldn't commit message:%v", err.Error())
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
