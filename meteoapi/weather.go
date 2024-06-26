package meteoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type WeatherResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	CurrentUnits         struct {
		Time          string `json:"time"`
		Temperature2M string `json:"temperature_2m"`
	} `json:"current_units"`
	Current struct {
		Time          string  `json:"time"`
		Temperature2M float64 `json:"temperature_2m"`
	} `json:"current"`
	Daily struct {
		Time             []string  `json:"time"`
		Temperature2MMax []float64 `json:"temperature_2m_max"`
		Temperature2MMin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

var (
	weatherRequestURL = `https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m&daily=temperature_2m_max,temperature_2m_min&temperature_unit=fahrenheit&wind_speed_unit=mph&precipitation_unit=inch`
)

func GetWeatherInfo(httpClient HttpClient, latitude, longitude float64) (*WeatherResponse, error) {
	var weatherResponse WeatherResponse

	url := fmt.Sprintf(weatherRequestURL, latitude, longitude)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Printf("request unsuccessful: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("unable to read response body: %v", err)
		return nil, err
	}

	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Printf("request unsuccessful: %v", err)
		return nil, err
	}
	return &weatherResponse, nil
}
