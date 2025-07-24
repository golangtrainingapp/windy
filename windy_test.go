package windy_test

import (
	"fmt"
	"github.com/golangtrainingapp/windy"
	"github.com/stretchr/testify/assert"
	_ "go/types"
	"net/http"
	"os"
	"testing"
)

const WINDYAPI_ENDPOINT = "https://api.windy.com/api/point-forecast/v2"

func TestBuildRequestReturnsRequestWithLatLongAndKey(t *testing.T) {
	t.Parallel()
	apiKey := ReturnApiKey()
	if apiKey == "" {
		t.Fatal("Api key not set")
	}

	request, err := windy.BuildRequest(53.1900, -112.2500, apiKey, "POST", WINDYAPI_ENDPOINT)
	if err != nil {
		fmt.Println(err)
		return
	}

	if request.Method != http.MethodPost {
		t.Errorf("Expected POST method, got %s", request.Method)
	}

	if request.URL.String() != "https://api.windy.com/api/point-forecast/v2" {
		t.Fatalf("WindyJsonResp does not match WINDYAPI_ENDPOINT")
	}

}

func TestValidateInputParametersFromRequest(t *testing.T) {
	t.Parallel()

	apikey := ReturnApiKey()
	if apikey == "" {
		t.Fatal("Missing api key")
	}

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
	apikey := ReturnApiKey()
	if apikey == "" {
		t.Fatal("Missing api key")
	}
	//Make a change in Config with invalid endpoint for example from v2 to v1
	endPoint := "https://api.windy.com/api/point-forecast/v1"
	req, _ := windy.BuildRequest(53.1900, -112.2500, apikey, "POST", endPoint)
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

}

func ReturnApiKey() string {
	var apiKey string
	if os.Getenv("WINDY_API_KEY") == "" {
		_ = os.Setenv("WINDY_API_KEY", "mxJW8fEadecqILVj7RWBdhUfJ38Ou0Bv")
	}
	apiKey = os.Getenv("WINDY_API_KEY")
	return apiKey
}
