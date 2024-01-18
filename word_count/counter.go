package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TTraveller7/invokerlib"
)

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	invokerlib.ProcessorHandle(w, r, counterProcess, counterInit)
}

func counterInit() {
	fmt.Printf("counterInit called")
}

func counterProcess(ctx context.Context, record *invokerlib.Record) error {
	fmt.Printf("%s\n", string(record.Value))
	return nil
}
