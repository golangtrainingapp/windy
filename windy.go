// The windy package is a simple http server package that makes request to windyapi.com
// and retrieves the real-time weather for a given latitude and longitude.
// When the client consumes this package they need to have latitude, longitude and apikey
package windy

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"strings"
)

// Windy Model
type Windy_Realtime_Report struct {
	Ts    []int64 `json:"ts"`
	Units struct {
		TempSurface             string `json:"temp-surface"`
		Temp1000H               string `json:"temp-1000h"`
		Temp800H                string `json:"temp-800h"`
		Temp400H                string `json:"temp-400h"`
		Temp200H                string `json:"temp-200h"`
		DewpointSurface         string `json:"dewpoint-surface"`
		Dewpoint1000H           string `json:"dewpoint-1000h"`
		Dewpoint800H            string `json:"dewpoint-800h"`
		Dewpoint400H            string `json:"dewpoint-400h"`
		Dewpoint200H            string `json:"dewpoint-200h"`
		Past3HprecipSurface     string `json:"past3hprecip-surface"`
		Past3HconvprecipSurface string `json:"past3hconvprecip-surface"`
		Past3HsnowprecipSurface string `json:"past3hsnowprecip-surface"`
		WindUSurface            string `json:"wind_u-surface"`
		WindU1000H              string `json:"wind_u-1000h"`
		WindU800H               string `json:"wind_u-800h"`
		WindU400H               string `json:"wind_u-400h"`
		WindU200H               string `json:"wind_u-200h"`
		WindVSurface            string `json:"wind_v-surface"`
		WindV1000H              string `json:"wind_v-1000h"`
		WindV800H               string `json:"wind_v-800h"`
		WindV400H               string `json:"wind_v-400h"`
		WindV200H               string `json:"wind_v-200h"`
		GustSurface             string `json:"gust-surface"`
		CapeSurface             string `json:"cape-surface"`
		PtypeSurface            any    `json:"ptype-surface"`
		LcloudsSurface          string `json:"lclouds-surface"`
		McloudsSurface          string `json:"mclouds-surface"`
		HcloudsSurface          string `json:"hclouds-surface"`
		RhSurface               string `json:"rh-surface"`
		Rh1000H                 string `json:"rh-1000h"`
		Rh800H                  string `json:"rh-800h"`
		Rh400H                  string `json:"rh-400h"`
		Rh200H                  string `json:"rh-200h"`
		GhSurface               string `json:"gh-surface"`
		Gh1000H                 string `json:"gh-1000h"`
		Gh800H                  string `json:"gh-800h"`
		Gh400H                  string `json:"gh-400h"`
		Gh200H                  string `json:"gh-200h"`
		PressureSurface         string `json:"pressure-surface"`
	} `json:"units"`
	TempSurface             []float64 `json:"temp-surface"`
	Temp1000H               []float64 `json:"temp-1000h"`
	Temp800H                []float64 `json:"temp-800h"`
	Temp400H                []float64 `json:"temp-400h"`
	Temp200H                []float64 `json:"temp-200h"`
	DewpointSurface         []float64 `json:"dewpoint-surface"`
	Dewpoint1000H           []float64 `json:"dewpoint-1000h"`
	Dewpoint800H            []float64 `json:"dewpoint-800h"`
	Dewpoint400H            []float64 `json:"dewpoint-400h"`
	Dewpoint200H            []float64 `json:"dewpoint-200h"`
	Past3HprecipSurface     []float64 `json:"past3hprecip-surface"`
	Past3HconvprecipSurface []float64 `json:"past3hconvprecip-surface"`
	Past3HsnowprecipSurface []float64 `json:"past3hsnowprecip-surface"`
	WindUSurface            []float64 `json:"wind_u-surface"`
	WindU1000H              []float64 `json:"wind_u-1000h"`
	WindU800H               []float64 `json:"wind_u-800h"`
	WindU400H               []float64 `json:"wind_u-400h"`
	WindU200H               []float64 `json:"wind_u-200h"`
	WindVSurface            []float64 `json:"wind_v-surface"`
	WindV1000H              []float64 `json:"wind_v-1000h"`
	WindV800H               []float64 `json:"wind_v-800h"`
	WindV400H               []float64 `json:"wind_v-400h"`
	WindV200H               []float64 `json:"wind_v-200h"`
	GustSurface             []float64 `json:"gust-surface"`
	CapeSurface             []float64 `json:"cape-surface"`
	PtypeSurface            []int     `json:"ptype-surface"`
	LcloudsSurface          []float64 `json:"lclouds-surface"`
	McloudsSurface          []float64 `json:"mclouds-surface"`
	HcloudsSurface          []float64 `json:"hclouds-surface"`
	RhSurface               []float64 `json:"rh-surface"`
	Rh1000H                 []float64 `json:"rh-1000h"`
	Rh800H                  []float64 `json:"rh-800h"`
	Rh400H                  []float64 `json:"rh-400h"`
	Rh200H                  []float64 `json:"rh-200h"`
	GhSurface               []float64 `json:"gh-surface"`
	Gh1000H                 []float64 `json:"gh-1000h"`
	Gh800H                  []float64 `json:"gh-800h"`
	Gh400H                  []float64 `json:"gh-400h"`
	Gh200H                  []float64 `json:"gh-200h"`
	PressureSurface         []float64 `json:"pressure-surface"`
	Warning                 string    `json:"warning"`
}

const WINDY_API_ENDPOINT = "https://api.windy.com/api/point-forecast/v2"

// The GetWeather function is responsible for retrieving the real-time weather response for a given
// latitude and longitude from WindyAPI.com. The apiKey is also necessary as the validation is
// performed by WindyAPI.com. An POST request is sent to WindyAPI.com to retrieve the real-time
// weather results
func GetWeather(latitude, longitude float64, apiKey string) (Windy_Realtime_Report, error) {

	if !math.IsNaN(latitude) && !math.IsInf(latitude, 0) && !isValidLatitude(latitude) {
		return Windy_Realtime_Report{}, errors.New("latitude must be a numeric value (between -90 and 90)")
	}
	if !math.IsNaN(longitude) && !math.IsInf(longitude, 0) && !isValidLongitude(longitude) {
		return Windy_Realtime_Report{}, errors.New("longitude must be a numeric value (between -180 and 180)")
	}
	if strings.Trim(apiKey, "") == "" {
		return Windy_Realtime_Report{}, errors.New("api key is required")
	}

	req, _ := BuildRequest(latitude, longitude, apiKey, "POST", WINDY_API_ENDPOINT)

	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Windy_Realtime_Report{}, err
	}
	respBytes, err := ParseWindyResponse(resp, err)
	if err != nil {
		return Windy_Realtime_Report{}, err
	}
	windyObj, err := UnMarshalResponseToWindyObject(respBytes)
	if err != nil {
		return Windy_Realtime_Report{}, err
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
func UnMarshalResponseToWindyObject(respBytes []byte) (Windy_Realtime_Report, error) {
	var resp Windy_Realtime_Report
	err := json.Unmarshal(respBytes, &resp)
	if err != nil {
		return Windy_Realtime_Report{}, err
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
