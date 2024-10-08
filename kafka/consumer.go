package kafka

import (
	"context"
	"encoding/json"
	"kafka-consumer/db"
	"kafka-consumer/models"
	"kafka-consumer/utils"
	"log"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/segmentio/kafka-go"
)

func ConsumeFromAllPartitions(topic string, bootstrapServer string, cache *map[string]models.Order, pool *pgxpool.Pool) {
	var mu sync.Mutex
	conn, err := kafka.DialLeader(context.Background(), "tcp", bootstrapServer, topic, 0)
	if err != nil {
		// log.Fatalf("Failed to dial Kafka broker for consumer: %v\n", err)
		utils.FatalError(err,"Failed to dial Kafka broker for consumer")
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		// log.Fatalf("Failed to read partitions for topic %s: %v\n", topic, err)
		utils.FatalError(err,"Failed to read partitions for topic %s")
	}
	for _, p := range partitions {
		go consumer(topic, bootstrapServer, cache, pool, &mu, p.ID)
	}
}

func consumer(topic string, bootstrapServer string, cache *map[string]models.Order, pool *pgxpool.Pool, mu *sync.Mutex, partition int) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{bootstrapServer},
		Topic:     topic,
		Partition: partition,
		MaxBytes:  10e6,
	})
	r.SetOffset(kafka.LastOffset)
	log.Println("Reading...")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		var order models.Order
		err = json.Unmarshal(m.Value, &order)
		if err != nil {
			// log.Println("Failed to parse JSON: Invalid JSON Format")
			utils.Warn("Failed to parse JSON: Invalid JSON Format")
			continue
		}

		mu.Lock()
		(*cache)[order.OrderUID] = order
		mu.Unlock()
		go db.SaveOrder(pool, &order)
	}

	if err := r.Close(); err != nil {
		// log.Fatal("failed to close reader:", err)
		utils.FatalError(err,"failed to close reader")
	}
}
