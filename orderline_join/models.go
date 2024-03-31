package main

import "time"

type Order struct {
	Wid       int64 `json:"wid"`
	Did       int64 `json:"did"`
	Oid       int64 `json:"oid"`
	CarrierId int64 `json:"carrier_id"`
}

type Orderline struct {
	Wid       int64     `json:"wid"`
	Did       int64     `json:"did"`
	Oid       int64     `json:"oid"`
	OlNumber  int64     `json:"ol_number"`
	Iid       int64     `json:"iid"`
	DeliveryD time.Time `json:"delivery_d"`
	Amount    float64   `json:"amount"`
	SupplyWid int64     `json:"supply_wid"`
	Quantity  int64     `json:"quantity"`
	DistInfo  string    `json:"dist_info"`
}

type FullOrderline struct {
	Wid       int64     `json:"wid"`
	Did       int64     `json:"did"`
	Oid       int64     `json:"oid"`
	CarrierId int64     `json:"carrier_id"`
	OlNumber  int64     `json:"ol_number"`
	Iid       int64     `json:"iid"`
	DeliveryD time.Time `json:"delivery_d"`
	Amount    float64   `json:"amount"`
	SupplyWid int64     `json:"supply_wid"`
	Quantity  int64     `json:"quantity"`
	DistInfo  string    `json:"dist_info"`
}
