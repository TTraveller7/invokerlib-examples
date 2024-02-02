package main

import (
	"context"
	"net/http"

	"github.com/TTraveller7/invokerlib"
)

var (
	counterPc = &invokerlib.ProcessorCallbacks{
		OnInit:  counterInit,
		Process: counterProcess,
	}
	stateStore invokerlib.StateStore
)

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	invokerlib.ProcessorHandle(w, r, counterPc)
}

func counterInit() error {
	var err error
	stateStore, err = invokerlib.NewFreeCacheStateStore()
	if err != nil {
		return err
	}
	invokerlib.AddStateStore("counter_state_store", stateStore)
	return nil
}

func counterProcess(ctx context.Context, record *invokerlib.Record) error {
	return stateStore.Put(ctx, record.Key, record.Value)
}
