package main

import (
	"context"
	"log"

	"github.com/muradrmagomedov/wbstestproject/consumer"
	"github.com/muradrmagomedov/wbstestproject/producer"
	"github.com/muradrmagomedov/wbstestproject/server"
	"github.com/muradrmagomedov/wbstestproject/storage/postgresql"
	redis "github.com/muradrmagomedov/wbstestproject/storage/redisDB"
)

func gracefulShutdown(reader *consumer.KafkaConsumer, writer *producer.KafkaProducer, postgresDB *postgresql.DB, redisDB *redis.DB, server *server.GinServer) {
	err := postgresDB.DB.Close()
	if err != nil {
		log.Printf("error while closing postgressql connection:%s\n", err.Error())
	} else {
		log.Println("postgresql connection is closed")
	}
	err = redisDB.RedisClient.Close()
	if err != nil {
		log.Printf("error while closing redis connection:%s\n", err.Error())
	} else {
		log.Println("redis connection is closed")
	}
	err = reader.Conn.Close()
	if err != nil {
		log.Printf("error while closing kafka reader connection:%s\n", err.Error())
	} else {
		log.Println("kafka.reader connection is closed")
	}
	err = writer.Conn.Close()
	if err != nil {
		log.Printf("error while closing kafka writer connection:%s\n", err.Error())
	} else {
		log.Println("kafka.writer connection is closed")
	}
	err = server.Server.Shutdown(context.Background())
	if err != nil {
		log.Printf("error while closing server  connection:%s\n", err.Error())
	} else {
		log.Println("server connection is closed")
	}
}
