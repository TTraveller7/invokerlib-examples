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
	stateStore, err = invokerlib.NewRedisStateStore("state-redis")
	if err != nil {
		return err
	}
	invokerlib.AddStateStore("counter_state_store", stateStore)
	return nil
}

func counterProcess(ctx context.Context, record *invokerlib.Record) error {
	val, err := stateStore.Get(ctx, record.Key)
	var n uint64 = 0
	if err == nil {
		n = invokerlib.BytesToUint64(val)
	}
	newVal := invokerlib.Uint64ToBytes(n + 1)
	err = stateStore.Put(ctx, record.Key, newVal)
	return err
}
