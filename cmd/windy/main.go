package main

import (
	"encoding/json"
	"fmt"
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

	resp, err := windy.GetWeather(53.1900, -112.2500, apiKey)
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
	fmt.Println(resp)
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
