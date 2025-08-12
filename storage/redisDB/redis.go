package redisDB

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	entities "github.com/muradrmagomedov/wbstestproject"
	"github.com/muradrmagomedov/wbstestproject/storage/postgresqldb"
	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	postgresDB  *postgresqldb.PostgresqlDB
	redisClient *redis.Client
}

func NewRedis(postgresDB *postgresqldb.PostgresqlDB) *RedisDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}
	return &RedisDB{
		redisClient: rdb,
		postgresDB:  postgresDB,
	}
}

func (r *RedisDB) SetupDB() error {
	orders, err := r.postgresDB.GetOrders("15", context.Background())
	if err != nil {
		return err
	}
	for _, order := range orders {
		buf, err := json.Marshal(order)
		if err != nil {
			log.Printf("problem with writing to redis:%v\n", err.Error())
		}
		_, err = r.redisClient.Set(context.Background(), order.OrderUID, buf, 10*time.Minute).Result()
		if err != nil {
			log.Printf("problem with writing to redis:%v\n", err.Error())
		}
	}
	log.Print("loading to redis finished")
	return nil
}

func (r *RedisDB) GetOrderByUID(orderUID string, ctx context.Context) (entities.Order, error) {
	order := entities.Order{}
	val, err := r.redisClient.Get(ctx, orderUID).Result()
	if err != nil {
		return order, err
	}
	err = json.Unmarshal([]byte(val), &order)
	if err != nil {
		return order, fmt.Errorf("error while serialization order from redis:%v", err.Error())
	}
	log.Println("got order from redis")
	return order, nil
}
