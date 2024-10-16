package kafka

import (
	"context"
	"encoding/json"
	"kafka-consumer/models"
	"kafka-consumer/utils"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

func ProduceToAllPartitions(topic string, bootstrapServer string, waitMs int,msgAmount int) {
	var wg sync.WaitGroup
	conn, err := kafka.DialLeader(context.Background(), "tcp", bootstrapServer, topic, 0)
	if err != nil {
		utils.FatalError(err, "failed to dial leader")
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		log.Fatalf("Failed to read partitions for topic in producer %s: %v\n", topic, err)
	}

	for _, partition := range partitions {
		for i := 0; i < 1; i++ {
			wg.Add(1)
			go producer(topic, bootstrapServer, partition.ID, waitMs, &wg, msgAmount)
		}
	}
	wg.Wait()
}

func producer(topic string, bootstrapServer string, partition int, waitMs int, wg *sync.WaitGroup, amount int) {
	defer wg.Done()
	delivery := models.NewDelivery("Иванов Иван", "+9720000000", "12345", "Москва", "Улица", "МСК", "почта@почта.рус")
	payment := models.NewPayment("txn123", "req456", "RUB", "Платежная система да", 1000, time.Now(), "Банк Ивановых", 50, 950, 0)
	item1 := models.NewItem(12345, "track123", 100, "RID1", "Item1", 10, "M", 90, 98765, "Brand1", 200)
	item2 := models.NewItem(67890, "track456", 200, "RID2", "Item2", 20, "L", 180, 54321, "Brand2", 200)
	order := models.NewOrder("", "track789", "entry1", delivery, *payment, []models.Item{*item1, *item2}, "en_US", "signature", "cust123", "FedEx", "shard1", "SM1", "shard2", time.Now())

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{bootstrapServer},
		Topic:   topic,
		Async:   true,
	})
	defer writer.Close()

	for i := 0; i < amount; i++ {
		order.OrderUID = utils.RandStringBytes(32)
		j, err := json.Marshal(order)
		if err != nil {
			utils.FatalError(err, "Conversion failed")
		}

		msg := kafka.Message{
			Partition: partition,
			Value:     []byte(j),
		}
		err = writer.WriteMessages(context.Background(), msg)
		if err != nil {
			utils.FatalError(err, "failed to write messages")
		}
		time.Sleep(time.Millisecond * time.Duration(waitMs))
	}

}
