package server

import (
	"math"
	"reflect"
	"testing"
)

// Consistency test of constructResponse.
func Test_constructResponse(t *testing.T) {
	apiResp := &apiResponse{
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
		}{Temp: 56.76},
	}

	expected := &Response{
		Lat:         apiResp.Coord.Lat,
		Lon:         apiResp.Coord.Lon,
		Condition:   apiResp.Weather[0].Main,
		Temperature: moderateTemp,
	}

	resp := constructResponse(apiResp)

	if !reflect.DeepEqual(expected, resp) {
		t.Errorf("Unexpected response.\nexpected: %+v\nreceived: %+v",
			expected, resp)
	}
}

// Unit test of getGeneralTemp.
func Test_getGeneralTemp(t *testing.T) {
	tests := []struct {
		temp     float64
		expected string
	}{
		{-20, coldTemp},
		{math.MinInt64, coldTemp},
		{20, coldTemp},
		{coldMaxTemp, coldTemp},
		{60, moderateTemp},
		{hotMinTemp, hotTemp},
		{math.MaxInt64, hotTemp},
	}

	for i, tt := range tests {
		generalTemp := getGeneralTemp(tt.temp)

		if tt.expected != generalTemp {
			t.Errorf("Unexpected general temperature for temp %f (%d)."+
				"\nexected: %s\nreceived: %s",
				tt.temp, i, tt.expected, generalTemp)
		}
	}
}
