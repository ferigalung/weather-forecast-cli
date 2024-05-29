package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// Weather is a response of forecast.json api from weatherapi.com
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		}
	}
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch float64 `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				}
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	q := "-7.19,111.92" // Bojonegoro
	apiKey := os.Getenv("API_KEY")

	// if user pass an argument after app command
	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	res, err := http.Get("https://api.weatherapi.com/v1/forecast.json?key=" + apiKey + "&q=" + q + "&aqi=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	color.Green("Current Weather")
	fmt.Printf("%s, %s: %.0fC, %s\n\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	color.Green("Forecast")
	for _, hour := range hours {
		date := time.Unix(int64(hour.TimeEpoch), 0)
		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s - %.0fC, %.0f%% CoR, %s\n", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)
		if hour.ChanceOfRain > 50 {
			color.Red(message)
		} else {
			fmt.Print(message)
		}
	}
}
