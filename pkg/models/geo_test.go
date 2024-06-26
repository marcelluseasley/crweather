package models

import (
	"testing"
)

func TestGeo_IsEmpty(t *testing.T) {
	type fields struct {
		Name      string
		Latitude  float64
		Longitude float64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "empty struct",
			fields: fields{
				Name:      "",
				Latitude:  0.0,
				Longitude: 0.0,
			},
			want: true,
		},
		{
			name: "non-empty name",
			fields: fields{
				Name:      "New York",
				Latitude:  0.0,
				Longitude: 0.0,
			},
			want: false,
		},
		{
			name: "non-empty latitude and longitude",
			fields: fields{
				Name:      "",
				Latitude:  40.712776,
				Longitude: -74.005974,
			},
			want: false,
		},
		{
			name: "fully populated struct",
			fields: fields{
				Name:      "Los Angeles",
				Latitude:  34.052235,
				Longitude: -118.243683,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Geo{
				Name:      tt.fields.Name,
				Latitude:  tt.fields.Latitude,
				Longitude: tt.fields.Longitude,
			}
			if got := g.IsEmpty(); got != tt.want {
				t.Errorf("Geo.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
