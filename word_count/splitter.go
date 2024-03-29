package main

import (
	"context"
	"net/http"

	"github.com/TTraveller7/invokerlib/pkg/api"
	"github.com/TTraveller7/invokerlib/pkg/core"
	"github.com/TTraveller7/invokerlib/pkg/models"
	"github.com/tiktoken-go/tokenizer"
)

var enc tokenizer.Codec

var splitterPc = &models.ProcessorCallbacks{
	OnInit:  splitterInit,
	Process: splitterProcess,
}

func SplitterHandler(w http.ResponseWriter, r *http.Request) {
	api.ProcessorHandle(w, r, splitterPc)
}

func splitterInit() error {
	var err error
	enc, err = tokenizer.Get(tokenizer.Cl100kBase)
	return err
}

func splitterProcess(ctx context.Context, record *models.Record) error {
	valStr := string(record.Value())
	_, tokens, err := enc.Encode(valStr)
	if err != nil {
		return err
	}

	for _, token := range tokens {
		tokenBytes := []byte(token)
		r := models.NewRecord(tokenBytes, tokenBytes)
		if err := core.PassToDefaultOutputTopic(ctx, r); err != nil {
			return err
		}
	}
	return nil
}
