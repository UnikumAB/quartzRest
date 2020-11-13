package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNanoTime_Scan(t *testing.T) {
	type args struct {
		src interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantErr    bool
	}{
		{name: "input type int", args: struct{ src interface{} }{src: 1601018690000}, wantErr: false, wantOutput: "2020-09-25T09:24:50+02:00"},
		{name: "input type int64", args: struct{ src interface{} }{src: int64(1601018690024)}, wantErr: false, wantOutput: "2020-09-25T09:24:50.024+02:00"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			i := &NanoTime{}
			if err := i.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			parsedTime, err := time.Parse(time.RFC3339Nano, tt.wantOutput)
			if err != nil {
				t.Errorf("Failed to parse wanted time: %v", err)
			}
			location, err := time.LoadLocation("Europe/Stockholm")
			if err != nil {
				t.Errorf("Failed to load Location: %v", err)
			}
			assert.Equal(t, parsedTime.In(location), i.Time.In(location))
		})
	}
}

func TestNanoTime_Value(t *testing.T) {
	type fields struct {
		Time time.Time
	}

	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		{
			name: "Converting back 1",
			fields: struct{ Time time.Time }{
				Time: time.Date(2020, 10, 01, 21, 34, 22, 123,
					time.FixedZone("Europe/Stockholm", 2))},
			want:    1601588060123,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			i := NanoTime{
				Time: tt.fields.Time,
			}
			got, err := i.Value()

			if (err != nil) != tt.wantErr {
				t.Errorf("Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
