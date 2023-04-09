package main

import (
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
