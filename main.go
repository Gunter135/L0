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

func main() {
	utils.InitLogger()
	log.Println("Starting up...")
	log.Println("Reading config...")
	config, err := config.ReadConfig("./config/config.yaml")
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

	http.HandleFunc("/api/orders/{id}", handlers.OrderHandler(&cache))
	http.HandleFunc("/api/orders/produce/{amount}", handlers.ProduceHandler(config.KafkaConfig))

	utils.FatalError(http.ListenAndServe(":8080", nil), "Server ded(")
}
