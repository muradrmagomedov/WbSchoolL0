package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	entities "github.com/muradrmagomedov/wbstestproject"
	"github.com/muradrmagomedov/wbstestproject/storage/postgresql"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var cacheTime time.Duration
var ordersToFill string

type DB struct {
	postgresDB   *postgresql.DB
	RedisClient  *redis.Client
	OrdersToFill string
}

type Config struct {
	Addr     string
	Password string
	DB       int
}

func createConfig() Config {
	return Config{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.dbname"),
	}
}

func NewRedis(postgresDB *postgresql.DB) *DB {
	config := createConfig()
	setCacheTime()
	setOrderToFill()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}
	return &DB{
		RedisClient:  rdb,
		postgresDB:   postgresDB,
		OrdersToFill: ordersToFill,
	}
}

func (r *DB) SetupDB() error {
	orders, err := r.postgresDB.GetOrders(ordersToFill, context.Background())
	if err != nil {
		return err
	}
	for _, order := range orders {
		buf, err := json.Marshal(order)
		if err != nil {
			log.Printf("problem with writing to redis:%v\n", err.Error())
		}
		_, err = r.RedisClient.Set(context.Background(), order.OrderUID, buf, cacheTime).Result()
		if err != nil {
			log.Printf("problem with writing to redis:%v\n", err.Error())
		}
	}
	log.Print("loading to redis finished")
	return nil
}

func (r *DB) GetOrderByUID(orderUID string, ctx context.Context) (entities.Order, error) {
	order := entities.Order{}
	val, err := r.RedisClient.Get(ctx, orderUID).Result()
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

func (r *DB) AddOrder(order entities.Order, ctx context.Context) error {
	buf, err := json.Marshal(order)
	if err != nil {
		return err
	}
	_, err = r.RedisClient.Set(context.Background(), order.OrderUID, buf, cacheTime).Result()
	if err != nil {
		return err
	}
	return nil
}

func setCacheTime() {
	cT := viper.GetInt("redis.cache_time_in_min")
	cacheTime = time.Minute * time.Duration(cT)
}

func setOrderToFill() {
	ordersToFill = viper.GetString("redis.orders_to_fill")
}
