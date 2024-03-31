package main

import (
	"reflect"
	"testing"

	"github.com/TTraveller7/invokerlib/pkg/utils"
)

func Test_parseOrderline(t *testing.T) {
	type args struct {
		val []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Orderline
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				val: []byte("1,1,1,1,10859,2022-08-02 17:16:13.949,184,1,5,tsbfqsgkpnuvxyegeuvdgbt"),
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOrderline(tt.args.val)
			t.Log(utils.SafeJsonIndent(got))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOrderline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOrderline() = %v, want %v", got, tt.want)
			}
		})
	}
}
