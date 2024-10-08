package models

import "time"

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string    `json:"transaction"`
	RequestID    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       int       `json:"amount"`
	PaymentDT    time.Time `json:"payment_dt"`
	Bank         string    `json:"bank"`
	DeliveryCost int       `json:"delivery_cost"`
	GoodsTotal   int       `json:"goods_total"`
	CustomFee    int       `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              string    `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func NewOrder(orderUID, trackNumber, entry string,
	delivery Delivery, payment Payment,
	items []Item, locale, internalSignature,
	customerID, deliveryService, shardkey,
	smID, oofShard string, dateCreated time.Time) *Order {
	return &Order{
		OrderUID:          orderUID,
		TrackNumber:       trackNumber,
		Entry:             entry,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            locale,
		InternalSignature: internalSignature,
		CustomerID:        customerID,
		DeliveryService:   deliveryService,
		Shardkey:          shardkey,
		SmID:              smID,
		DateCreated:       dateCreated,
		OofShard:          oofShard,
	}
}

func NewDelivery(name, phone, zip, city, address, region, email string) Delivery {
	return Delivery{
		Name:    name,
		Phone:   phone,
		Zip:     zip,
		City:    city,
		Address: address,
		Region:  region,
		Email:   email,
	}
}

func NewPayment(transaction, requestID, currency, provider string, amount int,
	paymentDT time.Time, bank string, deliveryCost, goodsTotal, customFee int) *Payment {
	return &Payment{
		Transaction:  transaction,
		RequestID:    requestID,
		Currency:     currency,
		Provider:     provider,
		Amount:       amount,
		PaymentDT:    paymentDT,
		Bank:         bank,
		DeliveryCost: deliveryCost,
		GoodsTotal:   goodsTotal,
		CustomFee:    customFee,
	}
}

func NewItem(chrtID int, trackNumber string, price int, rid, name string, sale int,
	size string, totalPrice int, nmID int, brand string, status int) *Item {
	return &Item{
		ChrtID:      chrtID,
		TrackNumber: trackNumber,
		Price:       price,
		RID:         rid,
		Name:        name,
		Sale:        sale,
		Size:        size,
		TotalPrice:  totalPrice,
		NmID:        nmID,
		Brand:       brand,
		Status:      status,
	}
}
