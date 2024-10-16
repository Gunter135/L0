package handlers

import (
	"encoding/json"
	"fmt"
	"kafka-consumer/config"
	"kafka-consumer/kafka"
	"kafka-consumer/models"
	"kafka-consumer/utils"
	"net/http"
	"strconv"
)

func OrderHandler(cache *map[string]models.Order) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		orderUID, err := utils.GetOrdersParamFromURL(req.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		order, exists := (*cache)[orderUID]
		if exists {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(order)
		} else {
			fmt.Fprintf(w, "Couldn't find order with OrderId: %s", orderUID)
		}
	}
}
func ProduceHandler(config config.KafkaConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := utils.GetOrdersProduceParamFromURL(req.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		msgAmount, err := strconv.Atoi(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		go kafka.ProduceToAllPartitions(
			config.Topic,
			config.BootstrapServer,
			200,
			msgAmount,
		)
	}
}
