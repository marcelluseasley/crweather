package meteoapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type GeoResponse struct {
	Results []Results `json:"results"`
}

type Results struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	CountryCode string  `json:"country_code"`
	Country     string  `json:"country"`
}

var (
	geoRequestURL = `https://geocoding-api.open-meteo.com/v1/search?name=%s&count=10&language=en&format=json`
)

func GetGeoInfo(httpClient HttpClient, name string) (*GeoResponse, error) {
	var geoResponse GeoResponse
	resp, err := httpClient.Get(fmt.Sprintf(geoRequestURL, name))
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
	err = json.Unmarshal(body, &geoResponse)
	if err != nil {
		log.Printf("request unsuccessful: %v", err)
		return nil, err
	}
	return &geoResponse, nil
}
