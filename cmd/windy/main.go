package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/golangtrainingapp/windy"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
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
	PrintToConsoleUsingTabWriter(resp, latitude, longitude)
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

func PrintToConsoleUsingTabWriter(response windy.Windy_Realtime_Report, latitude, longitude float64) {
	color.Cyan("Latitude: %v, Longitude: %v", latitude, longitude)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	_, _ = fmt.Fprintln(w, "Timestamp\t Air Temp (C)\tDewpoint (C)\tWind (C)\tWind Gust (mph)\tCape (J/kg)\tAir Pressure (mb)\tPrecipitation\tHumidity (%)\tLow CLoud\tHigh Cloud\tMedium Cloud\tOverall Snow\t")
	count := 0
	for i := 0; i < len(response.Ts); i++ {
		t := time.UnixMilli(response.Ts[i])
		_, _ = fmt.Fprintln(w, t, "\t",
			fmt.Sprintf("%.2f", response.TempSurface[i]-273.15), "\t",
			fmt.Sprintf("%.2f", response.DewpointSurface[i]-273.15), "\t",
			fmt.Sprintf("%.2f", response.WindUSurface[i]), "\t",
			fmt.Sprintf("%.2f", response.GustSurface[i]), "\t",
			fmt.Sprintf("%.2f", response.CapeSurface[i]), "\t",
			fmt.Sprintf("%.2f", response.PressureSurface[i]/100), "\t",
			PrecipitationType(response.PtypeSurface[i]), "\t",
			fmt.Sprintf("%.2f", response.RhSurface[i]), "\t",
			fmt.Sprintf("%.2f", response.HcloudsSurface[i]), "\t",
			fmt.Sprintf("%.2f", response.McloudsSurface[i]), "\t",
			fmt.Sprintf("%.2f", response.LcloudsSurface[i]), "\t",
			fmt.Sprintf("%.2f", response.Past3HprecipSurface[i]))
		count++
	}
	fmt.Println("Total Weather Records: ", count)
	_ = w.Flush()
}

func PrecipitationType(iprecip int) string {
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
	return precipitation
}
