package main

import (
	"encoding/json"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/golangtrainingapp/windyv1/windy"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	config, err := LoadConfig("windy/windyclient.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	}

	apiKey := config.ServerInfo.Apikey
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

type Config struct {
	ServerInfo struct {
		Apikey string `json:"apikey"`
	}
}

func LoadConfig(yamlFile string) (*Config, error) {
	configFilePath, err := xdg.ConfigFile(yamlFile)
	if err != nil {
		return nil, err
	}
	yamlWindy, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return nil, err
	}

	if len(yamlWindy) == 0 {
		fmt.Printf("YAML file is empty\n")
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(yamlWindy, &cfg)
	if err != nil {
		fmt.Printf("Error unmarshalling YAML: %v\n", err)
		return nil, err
	}
	return &cfg, nil
}
