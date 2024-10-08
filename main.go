package main

import (
	"kafka-consumer/config"
	"kafka-consumer/db"
	"kafka-consumer/handlers"
	"kafka-consumer/kafka"
	"kafka-consumer/utils"
	"log"
	"net/http"
)

// Если успею:
// Автотесты, мок на хттп хендлер,кафку,бд
// Валидация json сообщения с кафки
func main() {
	utils.SetDefaultLogger()
	log.Println("Starting up...")
	log.Println("Reading config...")
	config, err := config.ReadConfig("config/config.yaml")
	if err != nil {
		utils.FatalError(err, "Error, couldn't read config file")
	}
	log.Println("Initializing Database")
	db.InitDB(config.PostgreSQLConfig)
	log.Println("Creating connection pool")
	pool, err := db.GetDBConnection(config.PostgreSQLConfig)
	if err != nil {
		utils.FatalError(err, "Error, couldn't initialize database connection pool")
	}
	cache, err := db.RetrieveCache(pool)
	if err != nil {
		// utils.FatalError(err, "Кароч опять все сдохло")
		utils.FatalError(err, "Couldn't retrieve cache from Database")
	}
	log.Printf("Cache contains %d orders", len(cache))
	defer pool.Close()
	go kafka.ConsumeFromAllPartitions(
		config.KafkaConfig.Topic,
		config.KafkaConfig.BootstrapServer,
		&cache,
		pool,
	)
	// Продюсер, через него закидывал сообщения для проверки логики
	go kafka.ProduceToAllPartitions(
		config.KafkaConfig.Topic,
		config.KafkaConfig.BootstrapServer,
		200,
	)

	http.HandleFunc("/order", handlers.OrderHandler(&cache))

	utils.FatalError(http.ListenAndServe(":8080", nil), "Server ded(")
}
