package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Delivery struct {
	Name    string `json:"name" validate:"required,min=2"`
	Phone   string `json:"phone" validate:"required,e164"`
	Zip     string `json:"zip" validate:"required,len=5"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required,min=5"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string    `json:"transaction" validate:"required"`
	RequestID    string    `json:"request_id" validate:"required"`
	Currency     string    `json:"currency" validate:"required,len=3"`
	Provider     string    `json:"provider" validate:"required"`
	Amount       int       `json:"amount" validate:"required,gte=0"`
	PaymentDT    time.Time `json:"payment_dt" validate:"required"`
	Bank         string    `json:"bank" validate:"required"`
	DeliveryCost int       `json:"delivery_cost" validate:"required,gte=0"`
	GoodsTotal   int       `json:"goods_total" validate:"required,gte=0"`
	CustomFee    int       `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id" validate:"required,gte=1"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"required,gte=0"`
	RID         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required,min=2"`
	Sale        int    `json:"sale" validate:"gte=0,lte=100"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"required,gte=0"`
	NmID        int    `json:"nm_id" validate:"required,gte=1"`
	Brand       string `json:"brand" validate:"required,min=2"`
	Status      int    `json:"status" validate:"required,gte=100,lte=599"`
}

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,dive"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	Shardkey          string    `json:"shardkey" validate:"required"`
	SmID              string    `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`
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

func ValidateOrder(order Order) error {
	validate := validator.New()
	err := validate.Struct(order)
	if err != nil {
		return err
	}
	return nil
}
