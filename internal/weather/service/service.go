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

// FetchAndProcessWeather obtains weather data based on input parameters, which can be a location name or latitude and longitude coordinates.
// The process involves:
// 1. Resolving the location name to geographical coordinates.
// 2. Fetching weather information for the resolved coordinates or directly using provided coordinates if the location name is not given.
// 3. Processing and storing the fetched weather data.
// 4. Retrieving the final weather data from the repository for return.
// An error is returned if any step in the process fails.
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

	// think of geo table kind of like a cache...
	geoData, err := ws.repo.RetrieveGeo(latitude, longitude)
	if err != nil {
		log.Printf("error: %v", err)
		log.Println("pulling latest entry from weather")
	}

	var weatherData models.Weather

	// shouldn't be empty, but double check
	// also check if the data is older than an hour from now...if so, its no good
	if !geoData.IsEmpty() && !geoData.Expired() {
		weatherData, err = ws.repo.RetrieveWeatherFromGeo(int64(geoData.ID))
		if err != nil {
			return models.Weather{}, err
		}
	} else {

		weatherResponse, err := getWeatherInfo(latitude, longitude)
		if err != nil {
			return models.Weather{}, err
		}

		weatherLookupID = ws.processWeather(geoResponse, weatherResponse)

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

// processWeather takes care of handling both weather and geographical data that we get from meteoapi.
// We need two things to make it work: a GeoResponse, which holds all our geographical info, and a WeatherResponse, which has our weather details.
// Here's what it does step by step:
// 1. It checks if we've got a valid GeoResponse. If we do, it grabs the location's name, its latitude and longitude, and the date we made the request, and saves all that info into our database.
// 2. If we don't have any weather data (meaning our WeatherResponse is empty), we just stop right there and return -1. This tells us something went wrong because we can't process weather data if we don't have any.
// 3. If we managed to save our geographical info successfully, we link that data with its corresponding weather data using the database ID.
// 4. From our WeatherResponse, we pull out the current weather conditions, where exactly in the world this is happening, what the weather's going to be like for the coming week, including the highs and lows for each day.
// 5. All this detailed weather data then gets stored in the database.
// 6. If something goes wrong while we're trying to save this data, we make a note of the error but keep going.
// In the end, we return the database ID for the weather data we've just stored. If we didn't have any weather data to start with, we return -1.
func (ws *Service) processWeather(gr *meteoapi.GeoResponse, wr *meteoapi.WeatherResponse) int64 {
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

	// need to join the data for storage; need to convert to string before joining
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
