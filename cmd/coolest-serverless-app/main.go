package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/BobyMCbobs/sample-ko-monorepo/pkg/common"
)

var (
	defaultLatitude  = "-41.300293"
	defaultLongitude = "174.780304"

	appStatusPageNotFound     = "Page not found"
	appStatusHealthHealthy    = "Healthy"
	appStatusHealthNotHealthy = "Not healthy"
	appStatusInternalError    = "Internal error"
)

func GetEnvOrDefault(env, input string) string {
	fromEnv, exists := os.LookupEnv(env)
	if exists {
		return fromEnv
	}
	return input
}

func GetLatitude() string {
	return GetEnvOrDefault("LATITUDE", defaultLatitude)
}

func GetLongitude() string {
	return GetEnvOrDefault("LONGITUDE", defaultLongitude)
}

type handlers struct {
	weatherMetrics *WeatherMetrics
}

func (h *handlers) pageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, appStatusPageNotFound)
}

func (h *handlers) getHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.weatherMetrics.Temperature == nil || h.weatherMetrics.WindSpeed == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, appStatusHealthNotHealthy)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, appStatusHealthHealthy)
}

func (h *handlers) getWeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.weatherMetrics.Temperature == nil || h.weatherMetrics.WindSpeed == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	body, err := json.Marshal(h.weatherMetrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, appStatusInternalError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(body))
}

type OpenMeteoResultCurrentWeather struct {
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windspeed"`
}

type OpenMeteoResult struct {
	CurrentWeather OpenMeteoResultCurrentWeather `json:"current_weather"`
}

type WeatherMetrics struct {
	Temperature *float64
	WindSpeed   *float64
}

type CoolestServerlessApp struct {
	server         *http.Server
	weatherMetrics *WeatherMetrics
	latitude       string
	longitude      string
}

func NewCoolestServerlessApp() *CoolestServerlessApp {
	weatherMetrics := &WeatherMetrics{}
	handlers := &handlers{weatherMetrics: weatherMetrics}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/_healthz", handlers.getHealth)
	mux.HandleFunc("/api/weather", handlers.getWeather)
	mux.HandleFunc("/{.*}", handlers.pageNotFound)

	handler := common.Logging(mux)
	return &CoolestServerlessApp{
		server: &http.Server{
			Addr:           ":8080",
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		weatherMetrics: weatherMetrics,
		latitude:       GetLatitude(),
		longitude:      GetLongitude(),
	}
}

func (c *CoolestServerlessApp) Run() {
	log.Printf("Listening on HTTP port '%v'\n", c.server.Addr)
	log.Fatal(c.server.ListenAndServe())
}

func (c *CoolestServerlessApp) updateWeatherMetrics() error {
	req, err := http.NewRequest(http.MethodGet, "https://api.open-meteo.com/v1/forecast", nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = url.Values{
		"latitude":        []string{c.latitude},
		"longitude":       []string{c.longitude},
		"current_weather": []string{"true"},
		"timezone":        []string{"Pacific/Auckland"},
	}.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed, %+v", string(respBody))
	}
	var result OpenMeteoResult
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return err
	}
	c.weatherMetrics.Temperature = &result.CurrentWeather.Temperature
	c.weatherMetrics.WindSpeed = &result.CurrentWeather.WindSpeed
	log.Printf("Weather updated: %+v %+v\n", *c.weatherMetrics.Temperature, *c.weatherMetrics.WindSpeed)
	return nil
}

func (c *CoolestServerlessApp) DoUpdateWeatherMetrics() {
	time.Sleep(time.Second * 10)
	for {
		if err := c.updateWeatherMetrics(); err != nil {
			log.Printf("Failed to update weather metrics, %v\n", err)
		}
		time.Sleep(time.Minute * 5)
	}
}

func main() {
	c := NewCoolestServerlessApp()
	go c.DoUpdateWeatherMetrics()
	c.Run()
}
