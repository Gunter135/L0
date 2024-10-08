package handlers

import (
	"encoding/json"
	"fmt"
	"kafka-consumer/models"
	"net/http"
)

func OrderHandler(cache *map[string]models.Order) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		orderUID := req.URL.Query().Get("id")
		order, exists := (*cache)[orderUID]
		if exists {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(order)
		} else {
			fmt.Fprintf(w, "Couldn't find order with OrderId: %s", orderUID)
		}
	}
}
