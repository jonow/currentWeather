package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	jww "github.com/spf13/jwalterweatherman"
)

const apiEndpointHost = "https://api.openweathermap.org/data/2.5/weather"

// StartServer starts the http server at the specified port.
func StartServer(port, apiKey string) error {
	http.HandleFunc("GET /weather/{coords}", getCurrentWeather(apiKey))

	jww.INFO.Printf("Starting server on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		// Expected error when server is instructed to shut down
		jww.INFO.Print(err)
	} else if err != nil {
		return err
	}
	return nil
}

// getCurrentWeather is the handler func for the current weather endpoint.
func getCurrentWeather(apiKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jww.DEBUG.Printf("Received request at %s", r.URL)
		lat, lon, err := parseCoordinates(r.PathValue("coords"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		apiResp, err := lookupCoordinates(lat, lon, apiEndpointHost, apiKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := constructResponse(apiResp)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			jww.ERROR.Printf("Failed to write response: %+v", err)
		}
	}
}

// lookupCoordinates looks up the weather for the coordinates using the
// specified API key. It returns the weather condition and the temperature in
// Kelvin.
func lookupCoordinates(lat, lon, host, apiKey string) (*apiResponse, error) {
	jww.DEBUG.Printf("Looking up %s, %s", lat, lon)

	// Construct endpoint
	endpoint := host + "?lat=" + lat + "&lon=" + lon + "&appid=" + apiKey +
		"&units=imperial"

	// Request weather from endpoint
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get from endpoint")
	}
	defer closeAndLogResponse(resp.Body)

	// Check OK status
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(
			parseResponseError(resp), "Failed to communicate with API endpoint")
	}

	var apiResp apiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to JSON decode response")
	}

	return &apiResp, nil
}

type apiResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

// parseCoordinates parses the comma seperated latitude/longitude coordinates.
func parseCoordinates(value string) (lat, lon string, err error) {
	coords := strings.Split(value, ",")
	if len(coords) != 2 {
		return "", "", errors.New("malformed coordinates")
	}

	return strings.TrimSpace(coords[0]), strings.TrimSpace(coords[1]), nil
}

// parseResponseError takes a response reader and tries to parse a Failure.
func parseResponseError(resp *http.Response) error {
	errStruct := struct {
		Code    string `json:"cod"`
		Message string `json:"message"`
	}{}

	err := json.NewDecoder(resp.Body).Decode(&errStruct)
	if err != nil {
		return errors.Wrap(err, "Failed to read response")
	}

	return errors.Errorf("%s (code %s)", errStruct.Message, errStruct.Code)
}

// closeAndLogResponse closes the response body and prints an error to the log
// if it fails.
func closeAndLogResponse(rc io.ReadCloser) {
	err := rc.Close()
	if err != nil {
		jww.ERROR.Printf("Failed to close response: %+v", err)
	}
}
