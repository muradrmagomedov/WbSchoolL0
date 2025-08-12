package main

import (
	"log"
	"time"

	"github.com/muradrmagomedov/wbstestproject/consumer"
	"github.com/muradrmagomedov/wbstestproject/producer"
	"github.com/muradrmagomedov/wbstestproject/server"
	"github.com/muradrmagomedov/wbstestproject/storage/postgresqldb"
	"github.com/muradrmagomedov/wbstestproject/storage/redisDB"
)

const addr = "localhost:8001"

func main() {
	writer := producer.NewProducer()
	go func() {
		for {
			writer.KafkaDialler()
			time.Sleep(5 * time.Second)
		}
	}()
	reader := consumer.NewConsumer()
	postgresqldb := postgresqldb.NewPostgresqlDB()
	err := postgresqldb.Connection()
	if err != nil {
		log.Println(err)
	}
	go reader.ReadMessage(*postgresqldb)

	redisDB := redisDB.NewRedis(postgresqldb)
	err = redisDB.SetupDB()
	if err != nil {
		log.Println(err)
	}
	server := server.NewRouter(postgresqldb, redisDB)
	log.Println("Server starting...")
	server.Run(addr)
}
