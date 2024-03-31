package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/TTraveller7/invokerlib/pkg/api"
	"github.com/TTraveller7/invokerlib/pkg/core"
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
	olBytes, _ := sonic.Marshal(ol)
	orderlineId := fmt.Sprintf("%v_%v", ol.oid, ol.olNumber)
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
		return nil, fmt.Errorf("parse orderline failed: length of words less than 10")
	}
	ol := &Orderline{
		wid:       SafeParseInt64(words[0]),
		did:       SafeParseInt64(words[1]),
		oid:       SafeParseInt64(words[2]),
		olNumber:  SafeParseInt64(words[3]),
		iid:       SafeParseInt64(words[4]),
		deliveryD: SafeParseTime(words[5]),
		amount:    SafeParseFloat64(words[6]),
		supplyWid: SafeParseInt64(words[7]),
		quantity:  SafeParseInt64(words[8]),
		distInfo:  words[9],
	}
	return ol, nil
}
