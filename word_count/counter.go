package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TTraveller7/invokerlib"
)

var counterPc = &invokerlib.ProcessorCallbacks{
	OnInit:  counterInit,
	Process: counterProcess,
}

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	invokerlib.ProcessorHandle(w, r, counterPc)
}

func counterInit() {
	fmt.Printf("counterInit called")
}

func counterProcess(ctx context.Context, record *invokerlib.Record) error {
	fmt.Printf("%s\n", string(record.Value))
	return nil
}
