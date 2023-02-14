package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestUnmarshalQuery(t *testing.T) {
	f, err := os.Open("test/query.json")
	if err != nil {
		panic(err)
	}
	type args struct {
		f io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantType reflect.Type
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			"basic marshalling test",
			args{f},
			reflect.TypeOf(query{}),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalQuery(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", got)
			gott := reflect.TypeOf(got)
			if !reflect.DeepEqual(gott, tt.wantType) {
				t.Errorf("UnmarshalQuery() = %v, want %v", gott, tt.wantType)
			}
		})
	}
}

func Test_query_MarshalQuery(t *testing.T) {
	type fields struct {
		Model     string
		Prompt    string
		Temp      float64
		MaxTokens int
	}
	tests := []struct {
		name    string
		fields  fields
		want    *bytes.Buffer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"serialization test",
			fields{
				"foo",
				"bar",
				0.5,
				100,
			},
			bytes.NewBuffer([]byte{}),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query{
				Model:     tt.fields.Model,
				Prompt:    tt.fields.Prompt,
				Temp:      tt.fields.Temp,
				MaxTokens: tt.fields.MaxTokens,
			}
			got, err := q.MarshalQuery()
			if (err != nil) != tt.wantErr {
				t.Errorf("query.MarshalQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", got)
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
