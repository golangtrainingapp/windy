// The windy package is a simple http server package that makes request to windyapi.com
// and retrieves the real-time weather for a given latitude and longitude.
// When the client consumes this package they need to have latitude, longitude and apikey
package windy

import (
	"encoding/json"
	"errors"
	Config "github.com/golangtrainingapp/windyv1/windy/Config"
	"github.com/golangtrainingapp/windyv1/windy/model"
	"io"
	"math"
	"net/http"
	"strings"
)

// The GetWeather function is responsible for retrieving the real-time weather response for a given
// latitude and longitude from WindyAPI.com. The apiKey is also necessary as the validation is
// performed by WindyAPI.com. An POST request is sent to WindyAPI.com to retrieve the real-time
// weather results
func GetWeather(latitude, longitude float64, apiKey string) (model.Windy_Realtime_Report, error) {
	config, err := Config.LoadConfig("windy/windy.yaml")
	if err != nil {
		return model.Windy_Realtime_Report{}, errors.New("unable to process the request, please contact the application support team")
	}
	endPoint := config.ServerInfo.Endpoint

	if !math.IsNaN(latitude) && !math.IsInf(latitude, 0) && !isValidLatitude(latitude) {
		return model.Windy_Realtime_Report{}, errors.New("latitude must be a numeric value (between -90 and 90)")
	}
	if !math.IsNaN(longitude) && !math.IsInf(longitude, 0) && !isValidLongitude(longitude) {
		return model.Windy_Realtime_Report{}, errors.New("longitude must be a numeric value (between -180 and 180)")
	}
	if strings.Trim(apiKey, "") == "" {
		return model.Windy_Realtime_Report{}, errors.New("api key is required")
	}

	req, _ := BuildRequest(latitude, longitude, apiKey, "POST", endPoint)

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
func BuildRequest(latitude, longitude float64, apiKey string, methodType string, endPoint string) (*http.Request, error) {
	buildJsonReq := buildAPIRequest(latitude, longitude, apiKey)
	httpReq, err := http.NewRequest(methodType, endPoint, strings.NewReader(buildJsonReq))
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

// Performs the latitude validation
func isValidLatitude(lat float64) bool {
	return lat >= -90.0 && lat <= 90.0
}

// Performs the longitude validation
func isValidLongitude(lon float64) bool {
	return lon >= -180.0 && lon <= 180.0
}
