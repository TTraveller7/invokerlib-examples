package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/TTraveller7/invokerlib"
)

var splitterPc = &invokerlib.ProcessorCallbacks{
	OnInit:  splitterInit,
	Process: splitterProcess,
}

func SplitterHandler(w http.ResponseWriter, r *http.Request) {
	invokerlib.ProcessorHandle(w, r, splitterPc)
}

func splitterInit() error {
	fmt.Printf("splitterInit called")
	return nil
}

func splitterProcess(ctx context.Context, record *invokerlib.Record) error {
	valStr := string(record.Value)
	words := strings.Split(valStr, " ")
	for _, word := range words {
		r := &invokerlib.Record{
			Key:   word,
			Value: []byte(word),
		}
		if err := invokerlib.PassToDefaultOutputTopic(ctx, r); err != nil {
			return err
		}
	}
	return nil
}
