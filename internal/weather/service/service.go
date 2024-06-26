package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"crweather/internal/weather/db/sqlite3"
	"crweather/meteoapi"
	"crweather/pkg/models"
)

type WeatherServicer interface {
	ProcessWeather(*meteoapi.GeoResponse, *meteoapi.WeatherResponse) int64
	FetchAndProcessWeather(input, lat, long string) (models.Weather, error)
}

type Service struct {
	repo sqlite3.Storager
}

func NewWeatherService(storager sqlite3.Storager) WeatherServicer {
	return &Service{
		repo: storager,
	}
}

func (ws *Service) FetchAndProcessWeather(input, lat, long string) (models.Weather, error) {
	var latitude, longitude float64
	var geoResponse *meteoapi.GeoResponse
	var weatherLookupID int64
	var err error

	if input == "" && lat != "" && long != "" {
		latitude, longitude, err = getLatLong(lat, long)
		if err != nil {
			return models.Weather{}, err
		}
	} else {
		geoResponse, err = getGeoInfo(input)
		if err != nil {
			return models.Weather{}, err
		}
		if len(geoResponse.Results) < 1 {
			return models.Weather{}, errors.New("no geo results")
		}
		latitude = geoResponse.Results[0].Latitude
		longitude = geoResponse.Results[0].Longitude
	}

	weatherResponse, err := getWeatherInfo(latitude, longitude)
	if err != nil {
		return models.Weather{}, err
	}

	weatherLookupID = ws.ProcessWeather(geoResponse, weatherResponse)
	geoData, err := ws.repo.RetrieveGeo(latitude, longitude)
	if err != nil {
		log.Printf("error: %v", err)
		log.Println("pulling latest entry from weather")
	}

	var weatherData models.Weather
	if !geoData.IsEmpty() {

		weatherData, err = ws.repo.RetrieveWeatherFromGeo(int64(geoData.ID))
		if err != nil {
			return models.Weather{}, err
		}
	} else {
		weatherData, err = ws.repo.RetrieveWeather(weatherLookupID)
		if err != nil {
			return models.Weather{}, err
		}
	}

	weatherData.Daily.Time = strings.Split(weatherData.DailyDates, ",")
	weatherData.Daily.Temperature2MMax = stringSliceToFloatSlice(strings.Split(weatherData.DailyMin, ","))
	weatherData.Daily.Temperature2MMin = stringSliceToFloatSlice(strings.Split(weatherData.DailyMax, ","))
	return weatherData, nil
}

func stringSliceToFloatSlice(input []string) []float64 {
	output := make([]float64, len(input))

	for i, s := range input {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			f = 0.0
		}
		output[i] = f
	}
	return output
}

func getLatLong(lat, long string) (float64, float64, error) {

	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return 0, 0, errors.New("invalid latitude")
	}
	longitude, err := strconv.ParseFloat(long, 64)
	if err != nil {
		return 0, 0, errors.New("invalid longitude")
	}
	return latitude, longitude, nil
}

func (ws *Service) ProcessWeather(gr *meteoapi.GeoResponse, wr *meteoapi.WeatherResponse) int64 {
	var geoData models.Geo
	var geoID int64
	var weatherData models.Weather

	if gr != nil && len(gr.Results) > 0 {
		geoData.Name = gr.Results[0].Name
		geoData.Latitude = gr.Results[0].Latitude
		geoData.Longitude = gr.Results[0].Longitude
		geoData.RequestDate = time.Now()
		geoID = ws.repo.StoreGeo(geoData)
	}

	if wr == nil {
		return -1
	}
	if geoID > 0 {
		weatherData.GeoID = geoID
	}
	weatherData.Current = float32(wr.Current.Temperature2M)
	weatherData.Latitude = wr.Latitude
	weatherData.Longitude = wr.Longitude
	weatherData.DailyDates = strings.Join(wr.Daily.Time, ",")

	dailyMin := make([]string, len(wr.Daily.Temperature2MMin))
	for i, temp := range wr.Daily.Temperature2MMin {
		dailyMin[i] = fmt.Sprintf("%.2f", temp)
	}
	weatherData.DailyMin = strings.Join(dailyMin, ",")

	dailyMax := make([]string, len(wr.Daily.Temperature2MMax))
	for i, temp := range wr.Daily.Temperature2MMax {
		dailyMax[i] = fmt.Sprintf("%.2f", temp)
	}
	weatherData.DailyMax = strings.Join(dailyMax, ",")
	idInserted, err := ws.repo.StoreWeather(weatherData)
	if err != nil {
		log.Printf("unable to store weather: %v", err)
	}
	return idInserted
}

func getGeoInfo(input string) (*meteoapi.GeoResponse, error) {
	geoResponse, err := meteoapi.GetGeoInfo(input)
	if err != nil {
		return nil, err
	}
	return geoResponse, nil
}

func getWeatherInfo(latitude, longitude float64) (*meteoapi.WeatherResponse, error) {
	weatherResponse, err := meteoapi.GetWeatherInfo(latitude, longitude)
	if err != nil {
		return nil, err
	}
	return weatherResponse, nil
}
