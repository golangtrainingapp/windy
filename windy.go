// The windy package is a simple http server package that makes request to windyapi.com
// and retrieves the real-time weather for a given latitude and longitude.
// When the client consumes this package they need to have latitude, longitude and apikey
package windy

import (
	"encoding/json"
	"errors"
	"github.com/golangtrainingapp/windyv1/model"
	"io"
	"net/http"
	"strings"
)

// WINDYAPI_ENDPOINT: windyapi.com endpoint
const WINDYAPI_ENDPOINT = "https://api.windy.com/api/point-forecast/v2"

// The GetWeather function is responsible for retrieving the real-time weather response for a given
// latitude and longitude from WindyAPI.com. The apiKey is also necessary as the validation is
// performed by WindyAPI.com. An POST request is sent to WindyAPI.com to retrieve the real-time
// weather results
func GetWeather(latitude, longitude float64, apiKey string) (model.Windy_Realtime_Report, error) {
	req, err := BuildRequest(latitude, longitude, apiKey, "POST")
	if err != nil {
		return model.Windy_Realtime_Report{}, err
	}
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Windy_Realtime_Report{}, err
	}
	respBytes, err := ParseWindyResponse(resp, err)
	if err != nil {
		return model.Windy_Realtime_Report{}, err
	}
	windyObj, err := UnMarshalResponseToWindyObject(respBytes)
	if err != nil {
		return model.Windy_Realtime_Report{}, err
	}
	return windyObj, nil
}

// The buildAPIRequest function builds the POST request for the WindyAPI.com website
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

// The BuildRequest function internally invokes buildAPIRequest method for building the POST request
func BuildRequest(latitude, longitude float64, apiKey string, methodType string) (*http.Request, error) {
	buildJsonReq := buildAPIRequest(latitude, longitude, apiKey)
	httpReq, err := http.NewRequest(methodType, WINDYAPI_ENDPOINT, strings.NewReader(buildJsonReq))
	if err != nil {
		return nil, err
	}
	return httpReq, nil
}

// The ParseWindyResponse function parses the response from WindyAPI.com. If the parsing is
// successful it returns the response else it returns the error
func ParseWindyResponse(resp *http.Response, err error) ([]byte, error) {
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	windyJsonResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return windyJsonResp, nil

}

// The UnMarshalResponseToWindyObject function unmarshal's the WindyAPI.com response
// to a WindyAPI structure. If the operation is successful it returns the structure else it
// returns error
func UnMarshalResponseToWindyObject(respBytes []byte) (model.Windy_Realtime_Report, error) {
	var resp model.Windy_Realtime_Report
	err := json.Unmarshal(respBytes, &resp)
	if err != nil {
		return model.Windy_Realtime_Report{}, err
	}
	return resp, nil
}
