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

var apiKey string

func main() {
	err := godotenv.Load()
	if err == nil {
		apiKey = os.Getenv("API_KEY")
	}
	q := "-7.19,111.92" // Bojonegoro

	// if user pass an argument after app command
	index := 0
	if len(os.Args) >= 2 {
		if os.Args[1] == "besok" {
			index = 1
		} else {
			q = os.Args[1]
		}
	}

	res, err := http.Get("https://api.weatherapi.com/v1/forecast.json?days=2&key=" + apiKey + "&q=" + q + "&aqi=no")
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
	degreeSymbol := "\u00B0"
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[index].Hour

	color.Green("Current Weather")
	fmt.Printf("%s, %s: %.0f%sC, %s\n\n", location.Name, location.Country, current.TempC, degreeSymbol, current.Condition.Text)

	if index == 1 {
		color.Green("Tomorrow's forecast")
	} else {
		color.Green("Today's forecast")
	}

	for _, hour := range hours {
		date := time.Unix(int64(hour.TimeEpoch), 0)
		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s - %.0f%sC, %.0f%%, %s\n", date.Format("15:04"), hour.TempC, degreeSymbol, hour.ChanceOfRain, hour.Condition.Text)
		if hour.ChanceOfRain > 50 && hour.ChanceOfRain <= 90 {
			color.Yellow(message)
		} else if hour.ChanceOfRain > 90 {
			color.Red(message)
		} else {
			fmt.Print(message)
		}
	}
}
