package server

// Response is the payload sent to the caller containing the current weather.
type Response struct {
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Condition   string  `json:"condition"`
	Temperature string  `json:"temperature"`
}

const (
	// The possible temperature for Response.Temperature.
	coldTemp     = "cold"
	moderateTemp = "moderate"
	hotTemp      = "hot"

	// The maximum temperature for Response.Temperature to be set to coldTemp.
	coldMaxTemp = 50

	// The minimum temperature for Response.Temperature to be set to hotTemp.
	hotMinTemp = 80
)

// constructResponse constructs the response with the weather condition and
// temperature.
func constructResponse(apiResp *apiResponse) *Response {
	var condition string
	if len(apiResp.Weather) > 0 {
		condition = apiResp.Weather[0].Main
	}

	return &Response{
		Lat:         apiResp.Coord.Lat,
		Lon:         apiResp.Coord.Lon,
		Condition:   condition,
		Temperature: getGeneralTemp(apiResp.Main.Temp),
	}
}

// getGeneralTemp determines the general temperature from the precise
// temperature in Fahrenheit.
func getGeneralTemp(temp float64) string {
	if temp <= coldMaxTemp {
		return coldTemp
	} else if temp >= hotMinTemp {
		return hotTemp
	} else {
		return moderateTemp
	}
}
