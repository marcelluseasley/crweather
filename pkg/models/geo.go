package models

import "time"

type Geo struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	RequestDate time.Time `json:"access_date"`
}

func (g Geo) IsEmpty() bool {
	return g.ID ==0 && g.Name == "" && g.Latitude == 0 && g.Longitude == 0
}
