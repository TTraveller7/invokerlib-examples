package main

import (
	"context"
	"net/http"

	"github.com/TTraveller7/invokerlib/pkg/api"
	"github.com/TTraveller7/invokerlib/pkg/core"
	"github.com/TTraveller7/invokerlib/pkg/models"
	"github.com/TTraveller7/invokerlib/pkg/utils"
	"github.com/bytedance/sonic"
)

var (
	orderlineJoinPc = &models.ProcessorCallbacks{
		OnInit: orderlineInit,
		Join:   orderlineJoin,
	}
	metricsClient *utils.MetricsClient
)

func OrderlineJoinHandler(w http.ResponseWriter, r *http.Request) {
	api.ProcessorHandle(w, r, orderlineJoinPc)
}

func orderlineInit() error {
	metricsClient = core.MetricsClient()
	return nil
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
	if order.OrderId != orderline.OrderId {
		return nil
	}
	fullOrderline := &FullOrderline{
		Wid:       order.Wid,
		Did:       order.Did,
		Oid:       order.Oid,
		OrderId:   order.OrderId,
		CarrierId: order.CarrierId,
		OlNumber:  orderline.OlNumber,
		Iid:       orderline.Iid,
		DeliveryD: orderline.DeliveryD,
		Amount:    orderline.Amount,
		SupplyWid: orderline.SupplyWid,
		Quantity:  orderline.Quantity,
		DistInfo:  orderline.DistInfo,
	}
	fullOlBytes, _ := sonic.Marshal(fullOrderline)
	newRecord := models.NewRecord(rightRecord.Key(), fullOlBytes)
	if err := core.PassToDefaultOutputTopic(ctx, newRecord); err != nil {
		return err
	}
	metricsClient.EmitCounter("num_of_full_orderline", "number of full orderline record assembled", 1)
	return nil
}
