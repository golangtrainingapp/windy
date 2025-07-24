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
	color.Cyan("Latitude: %v, Longitude: %v", latitude, longitude)
	color.Blue("-------------------------------------------\n")
	color.Cyan("A sample partial data is displayed to the console:\n")
	color.Green("Surface Temperature:  %v", response.TempSurface[0])
	color.Green("Surface Dewpoint: %v", response.DewpointSurface[0])
	color.Green("Surface Wind: %v", response.WindUSurface[0])
	color.Green("Surface Gust: %v", response.GustSurface[0])
	color.Green("Surface Cape: %v", response.CapeSurface[0])
	color.Green("Surface Pressure: %v", response.PressureSurface[0])
	color.Green("Surface Ptype: %v", response.PtypeSurface[0])
	color.Green("Surface Lcloud: %v", response.LcloudsSurface[0])
	color.Green("Surface Hcloud %v", response.HcloudsSurface[0])
	color.Green("Surface Mcloud %v", response.McloudsSurface[0])
	color.Green("Surface Precipitation %v", response.Past3HprecipSurface[0])

}
