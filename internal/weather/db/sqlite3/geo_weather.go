package sqlite3

import (
	"database/sql"
	"log"
	"os"

	"crweather/pkg/models"

	_ "github.com/mattn/go-sqlite3"
)

type Storager interface {
	StoreGeo(models.Geo) int64
	StoreWeather(models.Weather) (int64, error)
	RetrieveGeo(float64, float64) (models.Geo, error)
	RetrieveWeather(int64) (models.Weather, error)
	RetrieveWeatherFromGeo(int64) (models.Weather, error)
}

type Store struct {
	db *sql.DB
}

func NewDBStore() (Storager, error) {

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "geo_weather.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) StoreGeo(geoData models.Geo) int64 {
	stmt, err := s.db.Prepare("INSERT INTO geo(name, latitude, longitude, access_date) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Printf("unable to prepare query - store geodata: %v", err)
		return 0
	}
	defer stmt.Close()

	result, err := stmt.Exec(geoData.Name, geoData.Latitude, geoData.Longitude, geoData.RequestDate)
	if err != nil {
		log.Printf("unable to execute statement - store geodata: %v", err)
		return 0
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("unable to get last inserted id: %v", err)
	}
	return id
}

func (s *Store) StoreWeather(weatherData models.Weather) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO weather(geo_id, current, daily_dates, daily_min, daily_max, latitude, longitude) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(weatherData.GeoID, weatherData.Current, weatherData.DailyDates, weatherData.DailyMin, weatherData.DailyMax, weatherData.Latitude, weatherData.Longitude)
	if err != nil {
		log.Printf("unable to execute statement - store weather data: %v", err)
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("unable to get last inserted id: %v", err)
	}
	return id, nil
}

func (s *Store) RetrieveWeather(id int64) (models.Weather, error) {
	var weather models.Weather

	row := s.db.QueryRow("SELECT geo_id, current, daily_dates, daily_min, daily_max, latitude, longitude FROM weather WHERE id = ?", id)

	err := row.Scan(&weather.GeoID, &weather.Current, &weather.DailyDates, &weather.DailyMin, &weather.DailyMax, &weather.Latitude, &weather.Longitude)
	if err != nil {
		return models.Weather{}, err
	}

	return weather, nil
}

func (s *Store) RetrieveWeatherFromGeo(geoID int64) (models.Weather, error) {
	var weather models.Weather

	row := s.db.QueryRow("SELECT geo_id, current, daily_dates, daily_min, daily_max, latitude, longitude FROM weather WHERE geo_id = ?", geoID)

	err := row.Scan(&weather.GeoID, &weather.Current, &weather.DailyDates, &weather.DailyMin, &weather.DailyMax, &weather.Latitude, &weather.Longitude)
	if err != nil {
		return models.Weather{}, err
	}

	return weather, nil
}

func (s *Store) RetrieveGeo(latitude, longitude float64) (models.Geo, error) {
	var geo models.Geo

	row := s.db.QueryRow("SELECT id, latitude, longitude, name FROM geo WHERE latitude = ? AND longitude = ? ORDER BY id DESC", latitude, longitude)

	err := row.Scan(&geo.ID, &geo.Latitude, &geo.Longitude, &geo.Name)
	if err != nil {
		return models.Geo{}, err
	}

	return geo, nil
}
