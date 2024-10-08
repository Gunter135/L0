package db

import (
	"context"
	"fmt"
	"kafka-consumer/config"
	"kafka-consumer/models"
	"kafka-consumer/utils"
	"log"
	"sync"

	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DataBaseFactory struct {
	pool *pgxpool.Pool
}

var (
	instance *DataBaseFactory
	once     sync.Once
)

func SaveOrder(conn *pgxpool.Pool, order *models.Order) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `
		INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), `
		INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}
	for _, item := range order.Items {
		_, err = conn.Exec(context.Background(), `
			INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDBConnection(dbconfig config.PostgreSQLConfig) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		url := fmt.Sprintf("%s://%s:%s@%s:%s/%s", dbconfig.DatabaseType, dbconfig.User, dbconfig.Password, dbconfig.Host, dbconfig.Port, dbconfig.Database)
		config, err := pgxpool.ParseConfig(url)
		if err != nil {
			log.Fatalf("Unable to parse connection string: %v\n", err)
		}

		config.MaxConns = dbconfig.MaxConnections
		config.MinConns = dbconfig.MinConnections

		pool, poolErr := pgxpool.ConnectConfig(context.Background(), config)
		if poolErr != nil {
			log.Fatalf("Unable to create connection pool: %v\n", poolErr)
		}

		pingErr := pool.Ping(context.Background())
		if pingErr != nil {
			log.Fatalf("Unable to ping database: %v\n", pingErr)
		}
		instance = &DataBaseFactory{
			pool: pool,
		}

		log.Println("Database connection pool established with MaxConns:", config.MaxConns)
	})
	if instance == nil {
		return nil, err
	}

	return instance.pool, nil
}

func InitDB(dbconfig config.PostgreSQLConfig) {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s", dbconfig.DatabaseType, dbconfig.User,
		dbconfig.Password, dbconfig.Host, dbconfig.Port, dbconfig.Database)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		utils.FatalError(err, "Unable to connect to database")
	}
	defer conn.Close(context.Background())
	data, err := os.ReadFile("./config/" + dbconfig.DbInit)
	if err != nil {
		utils.FatalError(err, "Unable to open dbinit file: ")
	}

	if _, err := conn.Exec(context.Background(), string(data)); err != nil {
		utils.FatalError(err, "Failed to execute dbinit script")
	}

	log.Println("Database initialized successfully.")
}

func RetrieveCache(conn *pgxpool.Pool) (map[string]models.Order, error) {
	cache := make(map[string]models.Order)
	query := `
		SELECT 
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, 
			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, 
			o.oof_shard,
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
			p.transaction, p.request_id, p.currency, p.provider, p.amount, 
			p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
			i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, 
			i.size, i.total_price, i.nm_id, i.brand, i.status
		FROM orders o
		JOIN deliveries d ON o.order_uid = d.order_uid
		JOIN payments p ON o.order_uid = p.order_uid
		JOIN items i ON o.order_uid = i.order_uid
	`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error querying orders: %v", err)
	}
	for rows.Next() {
		var orderUID, trackNumber, entry, locale, internalSignature, customerID, deliveryService, shardkey, smID, oofShard string
		var dateCreated time.Time
		var delivery models.Delivery
		var payment models.Payment
		var item models.Item

		err = rows.Scan(
			&orderUID, &trackNumber, &entry, &locale, &internalSignature, &customerID, &deliveryService, &shardkey, &smID, &dateCreated, &oofShard,
			&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
			&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %v", err)
		}
		if retrievedOrder, exists := cache[orderUID]; exists {
			retrievedOrder.Items = append(retrievedOrder.Items, item)
		} else {
			newOrder := models.Order{
				OrderUID:          orderUID,
				TrackNumber:       trackNumber,
				Entry:             entry,
				Locale:            locale,
				InternalSignature: internalSignature,
				CustomerID:        customerID,
				DeliveryService:   deliveryService,
				Shardkey:          shardkey,
				SmID:              smID,
				DateCreated:       dateCreated,
				OofShard:          oofShard,
				Delivery:          delivery,
				Payment:           payment,
				Items:             []models.Item{item},
			}
			retrievedOrder.Items = append(retrievedOrder.Items, item)
			cache[orderUID] = newOrder
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return cache, nil
}
