package postgresqldb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	entities "github.com/muradrmagomedov/wbstestproject"
)

const (
	host     = "localhost"
	port     = 5431
	user     = "root"
	password = "secret"
	dbname   = "root"
)

type PostgresqlDB struct {
	db *sql.DB
}

func NewPostgresqlDB() *PostgresqlDB {
	return &PostgresqlDB{}
}

func (p *PostgresqlDB) Connection() error {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("cant connect to postgres database: %v", err.Error())
	}
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("cant ping postgres database: %v", err.Error())
	}
	p.db = db
	return nil
}

func (p *PostgresqlDB) ParseMessage(order string) (entities.Order, error) {
	orderStruct := entities.Order{}
	err := json.Unmarshal([]byte(order), &orderStruct)
	return orderStruct, err
}

func (p *PostgresqlDB) AddOrder(order entities.Order, ctx context.Context) {
	queryDelivery := `Insert into delivery(order_uid,name,phone,zip,city,address,region,email) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);`
	queryPayment := `Insert into payment(transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee) Values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);`
	queryItem := `Insert into items(chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status) Values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);`
	queryOrder := `Insert into orders(order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,oof_shard) Values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("error while starting transaction:%v\n", err.Error())
		return
	}

	_, err = tx.ExecContext(ctx, queryOrder, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.OofShard)
	if err != nil {
		log.Printf("error while writing order:%v to db:%s\n", order.Payment.Transaction, err.Error())
		err = tx.Rollback()
		if err != nil {
			log.Printf("error while making rollback:%v\n", err.Error())
		}
		return
	}
	_, err = tx.ExecContext(ctx, queryDelivery, order.Delivery.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		log.Printf("error while writing delivery:%v to db:%s\n", order.Delivery.OrderUID, err.Error())
		err = tx.Rollback()
		if err != nil {
			log.Printf("error while making rollback:%v\n", err.Error())
		}
		return
	}
	_, err = tx.ExecContext(ctx, queryPayment, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		log.Printf("error while writing payment:%v to db:%s\n", order.Payment.Transaction, err.Error())
		err = tx.Rollback()
		if err != nil {
			log.Printf("error while making rollback:%v\n", err.Error())
		}
		return
	}
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, queryItem, item.ChartID, item.TrackNumber, item.Price, item.RId, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
		if err != nil {
			log.Printf("error while writing item:%v to db:%s\n", item.ChartID, err.Error())
			err = tx.Rollback()
			if err != nil {
				log.Printf("error while making rollback:%v\n", err.Error())
			}
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("error while commiting transactionto db:%s\n", err.Error())
	}
}

func (p *PostgresqlDB) GetOrderByUID(orderUID string, ctx context.Context) (entities.Order, error) {
	order := entities.Order{}
	query := `select order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard from orders where order_uid=$1;`
	err := p.db.QueryRowContext(ctx, query, orderUID).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		return order, fmt.Errorf("error while getting order from db:%v", err.Error())
	}
	delivery, err := p.GetDelivery(orderUID, ctx)
	if err != nil {
		return order, err
	}
	payment, err := p.GetPayment(orderUID, ctx)
	if err != nil {
		return order, err
	}
	items, err := p.GetItems(order.TrackNumber, ctx)
	if err != nil {
		return order, err
	}
	order.Delivery = delivery
	order.Payment = payment
	order.Items = items
	log.Println("got order from postgres")
	return order, nil
}

func (p *PostgresqlDB) GetDelivery(orderUID string, ctx context.Context) (entities.Delivery, error) {
	delivery := entities.Delivery{}
	query := `select name,phone,zip,city,address,region,email from delivery where order_uid=$1;`
	err := p.db.QueryRowContext(ctx, query, orderUID).Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil {
		return delivery, fmt.Errorf("error while getting delivery from db:%v", err.Error())
	}
	return delivery, nil
}

func (p *PostgresqlDB) GetPayment(orderUID string, ctx context.Context) (entities.Payment, error) {
	payment := entities.Payment{}
	query := `select transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee from payment where transaction=$1;`
	err := p.db.QueryRowContext(ctx, query, orderUID).Scan(&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)
	if err != nil {
		return payment, fmt.Errorf("error while getting payment from db:%v", err.Error())
	}
	return payment, nil
}

func (p *PostgresqlDB) GetItems(trackNumber string, ctx context.Context) ([]entities.Item, error) {
	items := []entities.Item{}
	query := `select chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status from items where track_number=$1;`
	rows, err := p.db.QueryContext(ctx, query, trackNumber)
	if err != nil {
		return nil, fmt.Errorf("error while getting items from db:%v", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var item entities.Item
		err := rows.Scan(&item.ChartID, &item.TrackNumber, &item.Price, &item.RId, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
		if err != nil {
			return nil, fmt.Errorf("error while getting items from db:%v", err.Error())
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while getting items from db:%v", err.Error())
	}
	return items, nil
}

func (p *PostgresqlDB) GetOrders(Limit string, ctx context.Context) ([]entities.Order, error) {
	orders := []entities.Order{}
	query := `select order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard from orders ORDER BY id DESC Limit $1;`

	rows, err := p.db.QueryContext(ctx, query, Limit)
	if err != nil {
		return nil, fmt.Errorf("error while getting orders from db:%v", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var order entities.Order
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
		if err != nil {
			return nil, fmt.Errorf("error while getting items from db:%v", err.Error())
		}
		delivery, err := p.GetDelivery(order.OrderUID, ctx)
		if err != nil {
			return nil, err
		}
		payment, err := p.GetPayment(order.OrderUID, ctx)
		if err != nil {
			return nil, err
		}
		items, err := p.GetItems(order.TrackNumber, ctx)
		if err != nil {
			return nil, err
		}
		order.Delivery = delivery
		order.Payment = payment
		order.Items = items
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error while getting items from db:%v", err.Error())
	}

	return orders, nil
}
