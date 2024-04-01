package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/TTraveller7/invokerlib/pkg/api"
	"github.com/TTraveller7/invokerlib/pkg/core"
	"github.com/TTraveller7/invokerlib/pkg/logs"
	"github.com/TTraveller7/invokerlib/pkg/models"
	"github.com/bytedance/sonic"
)

var (
	orderlineParsePc = &models.ProcessorCallbacks{
		Process: orderlineParseProcess,
	}
)

func OrderlineParseHandler(w http.ResponseWriter, r *http.Request) {
	api.ProcessorHandle(w, r, orderlineParsePc)
}

func orderlineParseProcess(ctx context.Context, record *models.Record) (orderlineProcessErr error) {
	ol, err := parseOrderline(record.Value())
	if err != nil {
		return err
	}

	if ol.OlNumber == 1 {
		o := &Order{
			Wid:       ol.Wid,
			Did:       ol.Did,
			Oid:       ol.Oid,
			CarrierId: int64(rand.Int()),
		}
		oBytes, _ := sonic.Marshal(o)
		orderRecord := models.NewRecord(strconv.FormatInt(o.Oid, 10), oBytes)
		if err := core.PassToOutputTopic(ctx, "orderparse", orderRecord); err != nil {
			return err
		}
	}

	olBytes, _ := sonic.Marshal(ol)
	orderlineId := fmt.Sprintf("%v_%v", ol.Oid, ol.OlNumber)
	newRecord := models.NewRecord(orderlineId, olBytes)
	if err := core.PassToDefaultOutputTopic(ctx, newRecord); err != nil {
		return err
	}

	return nil
}

// 1,1,1,1,10859,2022-08-02 17:16:13.949,184,1,5,tsbfqsgkpnuvxyegeuvdgbt
func parseOrderline(val []byte) (*Orderline, error) {
	words := strings.Split(string(val), ",")
	if len(words) < 10 {
		logs.Printf("parse orderline failed: length of words less than 10: %s", string(val))
		return &Orderline{}, nil
	}
	ol := &Orderline{
		Wid:       SafeParseInt64(words[0]),
		Did:       SafeParseInt64(words[1]),
		Oid:       SafeParseInt64(words[2]),
		OlNumber:  SafeParseInt64(words[3]),
		Iid:       SafeParseInt64(words[4]),
		DeliveryD: SafeParseTime(words[5]),
		Amount:    SafeParseFloat64(words[6]),
		SupplyWid: SafeParseInt64(words[7]),
		Quantity:  SafeParseInt64(words[8]),
		DistInfo:  words[9],
	}
	return ol, nil
}
