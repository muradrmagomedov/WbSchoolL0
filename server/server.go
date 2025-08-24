package server

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muradrmagomedov/wbstestproject/storage/postgresql"
	redis "github.com/muradrmagomedov/wbstestproject/storage/redisDB"
)

type GinServer struct {
	Router     *gin.Engine
	postgresdb *postgresql.DB
	redisdb    *redis.DB
	Server     *http.Server
}

func NewServer(addr string, db *postgresql.DB, redisdb *redis.DB) *GinServer {
	GinServer := &GinServer{
		postgresdb: db,
		redisdb:    redisdb,
		Router:     gin.Default(),
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: GinServer.Router,
	}
	GinServer.Server = srv
	return GinServer
}

func (s *GinServer) Run(addr string) error {
	s.Router.LoadHTMLGlob("static/*")

	s.Router.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "client.html", gin.H{
			"title": "Пример страницы",
		})
	})

	s.Router.POST("/order", s.getOrder)
	err := http.ListenAndServe(addr, s.Router)
	return err
}

func (s *GinServer) getOrder(c *gin.Context) {
	var orderID struct {
		OrderId string `json:"order_id"`
	}
	err := c.ShouldBindJSON(&orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	order, err := s.redisdb.GetOrderByUID(orderID.OrderId, context.Background())
	if err != nil {
		log.Println(err.Error())
		order, err = s.postgresdb.GetOrderByUID(orderID.OrderId, context.Background())
		if err != nil {
			log.Println("order_uid", orderID.OrderId)
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Заказ с таким orderID не найден", "errorMessage": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера", "errorMessage": err.Error()})
			}
			log.Println(err.Error())
			return
		}
		err = s.redisdb.AddOrder(order, context.Background())
		if err != nil {
			log.Printf("problem with writing to redis:%v\n", err.Error())
		}
	}
	c.JSON(http.StatusOK, order)
}
