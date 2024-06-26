package rest

import (
	"html/template"
	"net/http"
	"os"

	"crweather/internal/weather/service"

	"github.com/go-chi/chi/v5"
)

type WeatherHandler struct {
	service service.WeatherServicer
}

func NewWeatherHandler(service service.WeatherServicer) WeatherHandler {
	return WeatherHandler{
		service: service,
	}
}

func (handler WeatherHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("alive"))
}

func (handler WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	var err error
	input := chi.URLParam(r, "input")
	lat := chi.URLParam(r, "lat")
	long := chi.URLParam(r, "long")

	weatherResponse, err := handler.service.FetchAndProcessWeather(input, lat, long)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	
	}
	templatePath := os.Getenv("TMPL_PATH")
	tmplWeather, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmplWeather.Execute(w, weatherResponse)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}
