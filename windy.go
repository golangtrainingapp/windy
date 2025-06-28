package windy

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

const WINDYAPI_ENDPOINT = "https://api.windy.com/api/point-forecast/v2"

func GetWeather(latitude, longitude float64, apiKey string) (string, error) {
	req, err := BuildRequest(latitude, longitude, apiKey, "POST")
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	windyAPIResponse, err := ParseWindyResponse(resp, err)
	if err != nil {
		return "", err
	}
	return windyAPIResponse, nil
}

func buildAPIRequest(latitude, longitude float64, apiKey string) string {
	mapRequest := make(map[string]any)
	mapRequest["lat"] = latitude
	mapRequest["lon"] = longitude
	mapRequest["model"] = "gfs"
	mapRequest["parameters"] = []string{"temp", "dewpoint", "precip", "convPrecip", "snowPrecip", "wind", "windGust", "cape", "ptype", "lclouds", "mclouds", "hclouds", "rh", "gh", "pressure"}
	mapRequest["levels"] = []string{"surface", "1000h", "800h", "400h", "200h"}
	mapRequest["key"] = apiKey
	jsonRequest, _ := json.Marshal(mapRequest)
	return string(jsonRequest)
}

func BuildRequest(latitude, longitude float64, apiKey string, methodType string) (*http.Request, error) {
	buildJsonReq := buildAPIRequest(latitude, longitude, apiKey)
	httpReq, err := http.NewRequest(methodType, WINDYAPI_ENDPOINT, strings.NewReader(buildJsonReq))
	if err != nil {
		return nil, err
	}
	return httpReq, nil
}
func ParseWindyResponse(resp *http.Response, err error) (string, error) {
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	windyJsonResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(windyJsonResp), nil

}
