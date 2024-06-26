package main

import (
	"log"
	"net/http"

	"crweather/internal/weather/db/sqlite3"
	"crweather/internal/weather/rest"
	"crweather/internal/weather/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	weatherStorage, err := sqlite3.NewDBStore()
	if err != nil {
		log.Fatalf("could not initialize the database: %v", err)
	}
	httpClient := &http.Client{}
	weatherService := service.NewWeatherService(httpClient, weatherStorage)
	handler := rest.NewWeatherHandler(weatherService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handler.Ping)

	r.Route("/weather", func(r chi.Router) {
		r.Get("/{input}", handler.GetWeather)
		r.Get("/{lat}/{long}", handler.GetWeather)
	})

	log.Println("listening on port 3000")
	http.ListenAndServe(":3000", r)
}
