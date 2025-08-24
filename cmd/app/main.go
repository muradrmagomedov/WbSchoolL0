package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/muradrmagomedov/wbstestproject/consumer"
	"github.com/muradrmagomedov/wbstestproject/producer"
	"github.com/muradrmagomedov/wbstestproject/server"
	"github.com/muradrmagomedov/wbstestproject/storage/postgresql"
	redis "github.com/muradrmagomedov/wbstestproject/storage/redisDB"
	"github.com/spf13/viper"
)

func main() {
	LoadConfig()
	addr := viper.GetString("addr")

	writer := producer.NewProducer()
	go func() {
		for {
			writer.KafkaDialler()
			time.Sleep(5 * time.Second)
		}
	}()
	reader := consumer.NewConsumer()
	db := postgresql.NewPostgresqlDB()
	err := db.Connection()
	if err != nil {
		log.Println(err)
	}
	go reader.ReadMessage(*db)

	cahceDB := redis.NewRedis(db)
	err = cahceDB.SetupDB()
	if err != nil {
		log.Println(err)
	}
	server := server.NewServer(addr, db, cahceDB)
	log.Println("Server starting...")
	go func() {
		err := server.Run(addr)
		if err != nil {
			log.Fatalf("problems while starting server:%s", err.Error())
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Stopping application...")
	<-stop
	gracefulShutdown(reader, writer, db, cahceDB, server)
	log.Println("Application stopped")
}
