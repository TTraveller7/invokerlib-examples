package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TTraveller7/invokerlib/pkg/api"
	"github.com/TTraveller7/invokerlib/pkg/core"
	"github.com/TTraveller7/invokerlib/pkg/models"
	"github.com/bytedance/sonic"
)

var (
	orderlineParsePc = &models.ProcessorCallbacks{
		Process: orderlineParseProcess,
	}
	ErrIllFormat = fmt.Errorf("ill format record")
	orderIds     = sync.Map{}
	logs         = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

func OrderlineParseHandler(w http.ResponseWriter, r *http.Request) {
	api.ProcessorHandle(w, r, orderlineParsePc)
}

func orderlineParseProcess(ctx context.Context, record *models.Record) (orderlineProcessErr error) {
	ol, err := parseOrderline(record.Value())
	if err != nil {
		return nil
	}

	if _, exists := orderIds.Load(ol.OrderId); !exists {
		o := &Order{
			Wid:       ol.Wid,
			Did:       ol.Did,
			Oid:       ol.Oid,
			OrderId:   ol.OrderId,
			CarrierId: int64(rand.Int()),
		}
		oBytes, _ := sonic.Marshal(o)
		orderRecord := models.NewRecord(o.OrderId, oBytes)
		if err := core.PassToOutputTopic(ctx, "orderparse", orderRecord); err != nil {
			return err
		}
		orderIds.Store(ol.OrderId, true)
	}

	olBytes, _ := sonic.Marshal(ol)
	orderlineId := fmt.Sprintf("%v_%v", ol.OrderId, ol.OlNumber)
	newRecord := models.NewRecord(orderlineId, olBytes)
	if err := core.PassToDefaultOutputTopic(ctx, newRecord); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

// 1,1,1,1,10859,2022-08-02 17:16:13.949,184,1,5,tsbfqsgkpnuvxyegeuvdgbt
func parseOrderline(val []byte) (*Orderline, error) {
	words := strings.Split(string(val), ",")
	if len(words) < 10 {
		logs.Printf("parse orderline failed: length of words less than 10: %s", string(val))
		return nil, ErrIllFormat
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
	ol.OrderId = fmt.Sprintf("%v_%v_%v", ol.Wid, ol.Did, ol.Oid)
	return ol, nil
}
