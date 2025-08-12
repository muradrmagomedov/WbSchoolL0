package producer

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	topic     = "my-topic"
	partition = 0
)

type Producer struct{}

func NewProducer() *Producer {
	return &Producer{}
}

func (p Producer) KafkaDialler() {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("Fatal to create dialLeader", err)
	}
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
	conn.Close()
}
