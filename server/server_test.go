package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// Unit test of lookupCoordinates.
func Test_lookupCoordinates(t *testing.T) {
	expected := &apiResponse{
		Coord: struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}{
			Lon: -118.5386,
			Lat: 34.164,
		},
		Weather: []struct {
			Main string `json:"main"`
		}{{Main: "Clear"}},
		Main: struct {
			Temp float64 `json:"temp"`
		}{Temp: 286.91},
	}
	s := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			if query.Get("lat") == "" {
				t.Error("Empty lat")
			}
			if query.Get("lon") == "" {
				t.Error("Empty lon")
			}
			if query.Get("appid") == "" {
				t.Error("Empty appid")
			}

			_, _ = w.Write([]byte("{\"coord\":{\"lon\":-118.5386,\"lat\":34.164},\"weather\":[{\"id\":800,\"main\":\"Clear\",\"description\":\"clear sky\",\"icon\":\"01n\"}],\"base\":\"stations\",\"main\":{\"temp\":286.91,\"feels_like\":286.35,\"temp_min\":284.03,\"temp_max\":289.45,\"pressure\":1011,\"humidity\":77},\"visibility\":10000,\"wind\":{\"speed\":3.09,\"deg\":190},\"clouds\":{\"all\":0},\"dt\":1708137491,\"sys\":{\"type\":2,\"id\":2000061,\"country\":\"US\",\"sunrise\":1708094315,\"sunset\":1708133890},\"timezone\":-28800,\"id\":5395244,\"name\":\"Sherman Oaks\",\"cod\":200}"))
		}))
	defer s.Close()

	resp, err := lookupCoordinates("34.164", "-118.5386", s.URL, "apiKey")
	if err != nil {
		t.Errorf("Error looking up coordinates: %+v", err)
	}

	if !reflect.DeepEqual(expected, resp) {
		t.Errorf("Unexpected response.\nexpected: %+v\nreceived: %v",
			expected, resp)
	}
}

// Consistency test of parseCoordinates.
func Test_parseCoordinates(t *testing.T) {
	tests := []struct{ coords, lat, lon string }{
		{"39.8097343,-98.5556199", "39.8097343", "-98.5556199"},
		{"37.4144412, -108.0295765", "37.4144412", "-108.0295765"},
		{"-43, -103", "-43", "-103"},
	}

	for i, tt := range tests {
		lat, lon, err := parseCoordinates(tt.coords)
		if err != nil {
			t.Errorf("Failed to parse coordinates %q (%d)."+
				"\nexpected: %s, %s\nreceived: %s, %s",
				tt.coords, i, tt.lat, tt.lon, lat, lon)
		}
	}
}

// Error case: Tests that parseCoordinates returns an error for malformed
// coordinates.
func Test_parseCoordinates_Errors(t *testing.T) {
	tests := []string{
		"",
		"-108.0295765",
		"39.8097343,-98.5556199,37.4144412",
	}

	for i, coords := range tests {
		_, _, err := parseCoordinates(coords)
		if err == nil {
			t.Errorf("Nil error for invalid coordinates %q (%d).", coords, i)
		}
	}
}

// Consistency test of parseResponseError.
func Test_parseResponseError(t *testing.T) {
	expected := "Invalid date format (code 400)"
	data := `{"cod":"400", "message":"Invalid date format", "parameters": ["date"]}`

	err := parseResponseError(
		&http.Response{Body: io.NopCloser(strings.NewReader(data))})
	if err == nil || err.Error() != expected {
		t.Errorf("Unexpected error.\nexpected: %s\nreceived: %+v", expected, err)
	}
}
