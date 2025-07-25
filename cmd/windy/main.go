package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/golangtrainingapp/windy"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	var apiKey string
	if os.Getenv("WINDY_API_KEY") == "" {
		_ = os.Setenv("WINDY_API_KEY", "mxJW8fEadecqILVj7RWBdhUfJ38Ou0Bv")
	}
	apiKey = os.Getenv("WINDY_API_KEY")
	latitude := 53.1900
	longitude := -112.2500
	resp, err := windy.GetWeather(latitude, longitude, apiKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jsonPayload, _ := json.Marshal(resp)
	icao := "DCFG"
	path := filepath.Join(".", fmt.Sprintf("%s-%s.json", strconv.FormatInt(time.Now().Unix(), 10), icao))
	err = WriteToFile(path, jsonPayload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	PrintToConsole(resp, latitude, longitude)
	//fmt.Println(resp)
}

func WriteToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		println(err.Error())
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	_, _ = f.Write(data)
	return nil

}

func PrintToConsole(response windy.Windy_Realtime_Report, latitude, longitude float64) {
	// Use your own io.Writer output
	c := color.New(color.FgCyan).Add(color.Underline)

	fmt.Println()
	color.Cyan("Latitude: %v, Longitude: %v", latitude, longitude)
	color.Blue("-------------------------------------------\n")
	color.Cyan("A sample partial data is displayed to the console:\n")

	_, _ = c.Print("Air Temperature: ", fmt.Sprintf("%.2f", response.TempSurface[0]), " C,")
	_, _ = c.Print(" Dewpoint: ", fmt.Sprintf("%.2f", response.DewpointSurface[0]), " C,")
	_, _ = c.Print(" Wind: ", fmt.Sprintf("%.2f", response.WindUSurface[0]), " C,")
	_, _ = c.Print(" Wind Gust: ", fmt.Sprintf("%.2f", response.GustSurface[0]), " mph,")
	_, _ = c.Print(" Cape: ", fmt.Sprintf("%.2f", response.CapeSurface[0]), " J/Kg,")
	_, _ = c.Print(" Air Pressure: ", fmt.Sprintf("%.2f", response.PressureSurface[0]), " mb,")
	iprecip := response.PtypeSurface[0]
	var precipitation string
	switch iprecip {
	case 0:
		precipitation = "No Precipitation"
	case 1:
		precipitation = "Rain"
	case 3:
		precipitation = "Freeezing Rain"
	case 5:
		precipitation = "Snow"
	case 7:
		precipitation = "Mixture of Rain and snow"
	case 8:
		precipitation = "Ice Pellets"
	default:
		precipitation = "No Precipitation"
	}

	_, _ = c.Print(" Precipitation: ", precipitation, " ,")

	_, _ = c.Print(" Relative Humidity: ", fmt.Sprintf("%.2f", response.RhSurface[0]), " %,")
	_, _ = c.Print(" Low Cloud: ", fmt.Sprintf("%.2f", response.LcloudsSurface[0]), " ,")
	_, _ = c.Print(" High Cloud: ", fmt.Sprintf("%.2f", response.HcloudsSurface[0]), " ,")
	_, _ = c.Print(" Medium Cloud: ", fmt.Sprintf("%.2f", response.McloudsSurface[0]), " ,")
	_, _ = c.Print(" Overall snow for the preciding 3 hours: ", fmt.Sprintf("%.2f", response.Past3HprecipSurface[0]), " mm,")

}
