package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/TTraveller7/invokerlib/pkg/api"
	"github.com/TTraveller7/invokerlib/pkg/core"
	"github.com/TTraveller7/invokerlib/pkg/models"
	"github.com/bytedance/sonic"
)

var (
	orderParsePc = &models.ProcessorCallbacks{
		Process: orderParseProcess,
	}
)

func OrderParseHandler(w http.ResponseWriter, r *http.Request) {
	api.ProcessorHandle(w, r, orderParsePc)
}

func orderParseProcess(ctx context.Context, record *models.Record) (orderParseErr error) {
	o, err := parseOrder(record.Value())
	if err != nil {
		return err
	}
	oBytes, _ := sonic.Marshal(o)
	newRecord := models.NewRecord(strconv.FormatInt(o.oid, 10), oBytes)
	if err := core.PassToDefaultOutputTopic(ctx, newRecord); err != nil {
		return err
	}

	return nil
}

func parseOrder(val []byte) (*Order, error) {
	words := strings.Split(string(val), ",")
	if len(words) < 8 {
		return nil, fmt.Errorf("parse order failed: length of words less than 8")
	}
	o := &Order{
		wid:       SafeParseInt64(words[0]),
		did:       SafeParseInt64(words[1]),
		oid:       SafeParseInt64(words[2]),
		cid:       SafeParseInt64(words[3]),
		carrierId: SafeParseInt64(words[4]),
		olCnt:     SafeParseInt64(words[5]),
		allLocal:  SafeParseInt64(words[6]),
		entryD:    SafeParseTime(words[7]),
	}
	return o, nil
}
