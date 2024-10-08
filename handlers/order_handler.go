package handlers

import (
	"encoding/json"
	"fmt"
	"kafka-consumer/models"
	"kafka-consumer/utils"
	"net/http"
)

func OrderHandler(cache *map[string]models.Order) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		orderUID, err := utils.GetOrderIDFromURL(req.URL.Path)
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
