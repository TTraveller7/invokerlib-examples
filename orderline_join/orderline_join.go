package main

import (
	"context"
	"net/http"

	"github.com/TTraveller7/invokerlib/pkg/api"
	"github.com/TTraveller7/invokerlib/pkg/logs"
	"github.com/TTraveller7/invokerlib/pkg/models"
	"github.com/TTraveller7/invokerlib/pkg/utils"
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
	logs.Printf("orderline invoked: leftRecord: %s, rightRecord: %s", utils.SafeJsonIndent(leftRecord), utils.SafeJsonIndent(rightRecord))
	return nil
}
