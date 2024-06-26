package models

type Weather struct {
	ID         int64   `json:"id"`
	GeoID      int64   `json:"geo_id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Current    float32 `json:"current"`
	DailyDates string  `json:"daily_dates"`
	DailyMin   string  `json:"daily_min"`
	DailyMax   string  `json:"daily_max"`
	Daily      struct {
		Time             []string  `json:"time"`
		Temperature2MMax []float64 `json:"temperature_2m_max"`
		Temperature2MMin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}
