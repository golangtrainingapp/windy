package windy_test

import (
	"bytes"
	"fmt"
	"github.com/golangtrainingapp/windyv1"
	"io"
	"net/http"
	"testing"
)

const WINDYAPI_ENDPOINT = "https://api.windy.com/api/point-forecast/v2"

func TestBuildRequestReturnsRequestWithLatLongAndKey(t *testing.T) {
	request, err := windy.BuildRequest(53.1900, -112.2500, "mxJW8fEadecqILVj7RWBdhUfJ38Ou0Bv", "POST")
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
