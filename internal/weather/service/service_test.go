package service

import (
	"reflect"
	"testing"
)

func Test_stringSliceToFloatSlice(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "empty slice",
			args: args{input: []string{}},
			want: []float64{},
		},
		{
			name: "valid slice single element",
			args: args{input: []string{"1.5"}},
			want: []float64{1.5},
		},
		{
			name: "valid slice multiple elements",
			args: args{input: []string{"1.2", "3.4", "5.6"}},
			want: []float64{1.2, 3.4, 5.6},
		},
		{
			name: "invalid slice element",
			args: args{input: []string{"not_a_number"}},
			want: []float64{0},
		},
		{
			name: "mixed valid and invalid slice elements",
			args: args{input: []string{"3.14", "invalid", "2.71"}},
			want: []float64{3.14, 0, 2.71},
		},
		{
			name: "negative numbers",
			args: args{input: []string{"-1.23", "-4.56"}},
			want: []float64{-1.23, -4.56},
		},
		{
			name: "zero value",
			args: args{input: []string{"0"}},
			want: []float64{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringSliceToFloatSlice(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringSliceToFloatSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLatLong(t *testing.T) {
	type args struct {
		lat  string
		long string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		want1   float64
		wantErr bool
	}{
		{
			name: "valid coordinates",
			args: args{lat: "40.712776", long: "-74.005974"},
			want: 40.712776,
			want1: -74.005974,
			wantErr: false,
		},
		{
			name: "invalid latitude",
			args: args{lat: "not_a_latitude", long: "-74.005974"},
			want: 0,
			want1: 0,
			wantErr: true,
		},
		{
			name: "invalid longitude",
			args: args{lat: "40.712776", long: "not_a_longitude"},
			want: 0,
			want1: 0,
			wantErr: true,
		},
		{
			name: "both coordinates invalid",
			args: args{lat: "not_a_latitude", long: "not_a_longitude"},
			want: 0,
			want1: 0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getLatLong(tt.args.lat, tt.args.long)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLatLong() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getLatLong() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getLatLong() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
