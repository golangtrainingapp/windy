package windy_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golangtrainingapp/windyv1/windy"
	"github.com/golangtrainingapp/windyv1/windy/Config"
	_ "github.com/golangtrainingapp/windyv1/windy/model"
	"github.com/stretchr/testify/assert"
	_ "go/types"
	"io"
	"net/http"
	"testing"
)

const WINDYAPI_ENDPOINT = "https://api.windy.com/api/point-forecast/v2"

func TestInvalidConfigFile(t *testing.T) {
	t.Parallel()
	//Pass the invalid key pair in the Config file to simulate the error
	_, err := Config.LoadConfig("test/windy.yaml")
	if err != nil {
		assert.Error(t, errors.New("Unable to load the configuration file. Please contact the application support team."), err)
	}
}

func TestValidConfigFile(t *testing.T) {
	t.Parallel()

	//Pass the invalid key pair in the Config file to simulate the error
	config, err := Config.LoadConfig("windy/windy.yaml")
	if err != nil {
		assert.Error(t, errors.New("Unable to load the configuration file. Please contact the application support team."), err)
	}
	assert.YAMLEq(t, config.ServerInfo.Endpoint, "https://api.windy.com/api/point-forecast/v2")

}

func TestBuildRequestReturnsRequestWithLatLongAndKey(t *testing.T) {
	t.Parallel()
	request, err := windy.BuildRequest(53.1900, -112.2500, "mxJW8fEadecqILVj7RWBdhUfJ38Ou0Bv", "POST", WINDYAPI_ENDPOINT)
	if err != nil {
		fmt.Println(err)
		return
	}
	windyJsonResp, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatal(err)
	}

	sampleRequest := `{"key":"mxJW8fEadecqILVj7RWBdhUfJ38Ou0Bv","lat":53.19,"levels":["surface","1000h","800h","400h","200h"],"lon":-112.25,"model":"gfs","parameters":["temp","dewpoint","precip","convPrecip","snowPrecip","wind","windGust","cape","ptype","lclouds","mclouds","hclouds","rh","gh","pressure"]}`
	if !bytes.Equal(windyJsonResp, convertRequestToBytes(sampleRequest)) {
		t.Fatalf("WindyJsonResp does not match request body")
	}

	if request.Method != http.MethodPost {
		t.Errorf("Expected POST method, got %s", request.Method)
	}

	if request.URL.String() != "https://api.windy.com/api/point-forecast/v2" {
		t.Fatalf("WindyJsonResp does not match WINDYAPI_ENDPOINT")
	}

}

func convertRequestToBytes(req string) []byte {
	return []byte(req)
}

func TestValidateInputParametersFromRequest(t *testing.T) {
	t.Parallel()
	const apikey = "mxJW8fEadecqILVj7RWBdhUfJ38Ou0Bv"
	var requestTests = []struct {
		name          string
		latitude      float64
		longitude     float64
		apiKey        string
		validationMsg string
	}{
		{"Invalid latitude", 95, -112.2500, apikey, "latitude must be a numeric value (between -90 and 90)"},
		{"Invalid longitude", 53.1900, 200, apikey, "longitude must be a numeric value (between -180 and 180)"},
		{"Empty api key", 53.1900, -112.2500, "", "api key is required"},
		{"Invalid api key", 53.1900, -112.2500, "xJW8fEadecqILVj7RWBdhUfJ38Ou0Bv", "400 Bad Request"},
		{"Valid parameters", 53.1900, -112.2500, apikey, ""},
	}

	for _, tt := range requestTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := windy.GetWeather(tt.latitude, tt.longitude, tt.apiKey)
			if err != nil && !assert.Equal(t, tt.validationMsg, err.Error()) {
				t.Errorf("Expected %s, but got %s", tt.validationMsg, err.Error())
			}
		})

	}
}

func TestSimulateInvalidWindyEndPoint(t *testing.T) {
	t.Parallel()
	//Make a change in Config with invalid endpoint for example from v2 to v1
	endPoint := "https://api.windy.com/api/point-forecast/v1"
	req, _ := windy.BuildRequest(53.1900, -112.2500, "mxJW8fEadecqILVj7RWBdhUfJ38Ou0Bv", "POST", endPoint)
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

}
