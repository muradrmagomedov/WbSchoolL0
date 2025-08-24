package producer

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type KafkaProducer struct {
	Conn *kafka.Conn
}

func NewProducer() *KafkaProducer {
	return &KafkaProducer{}
}

type ProducerConfig struct {
	Addr      string
	Topic     string
	Partition int
}

func createConfig() ProducerConfig {
	return ProducerConfig{
		Addr:      viper.GetString("kafka.addr"),
		Topic:     viper.GetString("kafka.topic"),
		Partition: viper.GetInt("kafka.partition"),
	}
}

func (p *KafkaProducer) KafkaDialler() {
	config := createConfig()
	conn, err := kafka.DialLeader(context.Background(), "tcp", config.Addr, config.Topic, config.Partition)
	if err != nil {
		log.Fatal("Fatal to create dialLeader", err)
	}
	p.Conn = conn
	kfkMsg := kafka.Message{}
	if rand.Intn(10) > 3 {
		kfkMsg.Value = []byte(wrongMessage())
	} else {
		kfkMsg.Value = []byte(rightMessage())
	}
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal("Failed to add deadline", err)
	}
	_, err = conn.WriteMessages(kfkMsg)
	if err != nil {
		log.Fatal("Failed to write message to kafka", err)
	}
	err = conn.Close()
	if err != nil {
		log.Printf("problems when closing kafka dialer:%s", err.Error())
	}
}
