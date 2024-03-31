package main

import "time"

type Order struct {
	wid       int64
	did       int64
	oid       int64
	cid       int64
	carrierId int64
	olCnt     int64
	allLocal  int64
	entryD    time.Time
}

type Orderline struct {
	wid       int64
	did       int64
	oid       int64
	olNumber  int64
	iid       int64
	deliveryD time.Time
	amount    float64
	supplyWid int64
	quantity  int64
	distInfo  string
}

type FullOrderline struct {
	wid       int64
	did       int64
	oid       int64
	cid       int64
	carrierId int64
	olCnt     int64
	allLocal  int64
	entryD    time.Time
	olNumber  int64
	iid       int64
	deliveryD time.Time
	amount    float64
	supplyWid int64
	quantity  int64
	distInfo  string
}
