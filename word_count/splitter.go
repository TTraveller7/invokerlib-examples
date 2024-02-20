package main

import (
	"context"
	"net/http"

	"github.com/TTraveller7/invokerlib"
	"github.com/tiktoken-go/tokenizer"
)

var enc tokenizer.Codec

var splitterPc = &invokerlib.ProcessorCallbacks{
	OnInit:  splitterInit,
	Process: splitterProcess,
}

func SplitterHandler(w http.ResponseWriter, r *http.Request) {
	invokerlib.ProcessorHandle(w, r, splitterPc)
}

func splitterInit() error {
	var err error
	enc, err = tokenizer.Get(tokenizer.Cl100kBase)
	return err
}

func splitterProcess(ctx context.Context, record *invokerlib.Record) error {
	valStr := string(record.Value)
	_, tokens, err := enc.Encode(valStr)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		r := &invokerlib.Record{
			Key:   token,
			Value: []byte(token),
		}
		if err := invokerlib.PassToDefaultOutputTopic(ctx, r); err != nil {
			return err
		}
	}
	return nil
}
