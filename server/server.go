package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muradrmagomedov/wbstestproject/storage/postgresqldb"
	"github.com/muradrmagomedov/wbstestproject/storage/redisDB"
)

type Router struct {
	router     *gin.Engine
	postgresdb *postgresqldb.PostgresqlDB
	redisdb    *redisDB.RedisDB
}

func NewRouter(db *postgresqldb.PostgresqlDB, redisdb *redisDB.RedisDB) *Router {
	return &Router{
		postgresdb: db,
		redisdb:    redisdb,
		router:     gin.Default(),
	}
}

func (r *Router) Run(addr string) {
	r.router.LoadHTMLGlob("static/*")

	r.router.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "client.html", gin.H{
			"title": "Пример страницы",
		})
	})

	r.router.POST("/order", r.getOrder)
	r.router.Run(addr)
}

func (r *Router) getOrder(c *gin.Context) {
	var orderID struct {
		OrderId string `json:"order_id"`
	}
	err := c.ShouldBindJSON(&orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	order, err := r.redisdb.GetOrderByUID(orderID.OrderId, context.Background())
	if err != nil {
		order, err = r.postgresdb.GetOrderByUID(orderID.OrderId, context.Background())
		if err != nil {
			log.Println("order_uid", orderID.OrderId)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера", "errorMessage": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, order)
}
