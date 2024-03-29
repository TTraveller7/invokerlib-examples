package main

import (
	"context"
	"net/http"

	"github.com/TTraveller7/invokerlib/pkg/api"
	"github.com/TTraveller7/invokerlib/pkg/models"
	"github.com/TTraveller7/invokerlib/pkg/state"
	"github.com/TTraveller7/invokerlib/pkg/utils"
)

var (
	counterPc = &models.ProcessorCallbacks{
		OnInit:  counterInit,
		Process: counterProcess,
	}
	stateStore state.StateStore
)

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	api.ProcessorHandle(w, r, counterPc)
}

func counterInit() error {
	var err error
	stateStore, err = state.NewMemcachedStateStore("state-memcached")
	if err != nil {
		return err
	}
	state.AddStateStore("counter_state_store", stateStore)
	return nil
}

func counterProcess(ctx context.Context, record *models.Record) error {
	keyStr := string(record.Key())
	val, err := stateStore.Get(ctx, keyStr)
	var n uint64 = 0
	if err == nil {
		n = utils.BytesToUint64(val)
	}
	newVal := utils.Uint64ToBytes(n + 1)
	err = stateStore.Put(ctx, keyStr, newVal)
	return err
}
