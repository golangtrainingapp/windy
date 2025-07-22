package windy

import (
	"fmt"
	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	ServerInfo struct {
		Endpoint string `json:"endpoint"`
		ApiKey   string `json:"apikey"`
	}
}

func LoadConfig(yamlFile string) (*Config, error) {
	configFilePath, err := xdg.ConfigFile(yamlFile)
	if err != nil {
		return nil, err
	}
	yamlWindy, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	if len(yamlWindy) == 0 {
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
