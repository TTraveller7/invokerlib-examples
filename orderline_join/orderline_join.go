package main

import (
	"context"
	"net/http"

	"github.com/TTraveller7/invokerlib/pkg/api"
	"github.com/TTraveller7/invokerlib/pkg/logs"
	"github.com/TTraveller7/invokerlib/pkg/models"
	"github.com/TTraveller7/invokerlib/pkg/utils"
	"github.com/bytedance/sonic"
)

var (
	orderlineJoinPc = &models.ProcessorCallbacks{
		Join: orderlineJoin,
	}
)

func OrderlineJoinHandler(w http.ResponseWriter, r *http.Request) {
	api.ProcessorHandle(w, r, orderlineJoinPc)
}

func orderlineJoin(ctx context.Context, leftRecord *models.Record, rightRecord *models.Record) error {
	var order Order
	if err := sonic.Unmarshal(leftRecord.Value(), &order); err != nil {
		return err
	}
	var orderline Orderline
	if err := sonic.Unmarshal(rightRecord.Value(), &orderline); err != nil {
		return err
	}
	if order.Oid != orderline.Oid {
		return nil
	}
	fullOrderline := &FullOrderline{
		Wid:       order.Wid,
		Did:       order.Did,
		Oid:       order.Oid,
		CarrierId: order.CarrierId,
		OlNumber:  orderline.OlNumber,
		Iid:       orderline.Iid,
		DeliveryD: orderline.DeliveryD,
		Amount:    orderline.Amount,
		SupplyWid: orderline.SupplyWid,
		Quantity:  orderline.Quantity,
		DistInfo:  orderline.DistInfo,
	}
	logs.Printf("full order line produced: %s", utils.SafeJsonIndent(fullOrderline))
	return nil
}
