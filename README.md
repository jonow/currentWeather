# Current Weather

Current Weather is a simple HTTP server that returns the general weather
condition and temperature at an endpoint that takes in latitude and longitude
coordinates.

## Running Server

Start the server by specifying your OpenWeather API key using the `-k` flag.

```text
currentWeather -k <API key>
```

Other flags available are:

```text
Flags:
  -k, --apiKey string    OpenWeather API key.
  -h, --help             help for currentWeather
  -v, --logLevel int     Verbosity level for log printing (2+ = Trace, 1 = Debug, 0 = Info).
  -l, --logPath string   File path to save log file to. (default "-")
  -p, --port string      Port the server listens on. (default "9090")
```

## Accessing Endpoint

To get the current weather, connect to the following endpoint with the
coordinates seperated by a comma:

```text
localhost/weather/<lat coordinate>,<long coordinate>
```

Returns a JSON object containing the current weather condition (e.g., clear,
rain, snow) and the general temperature (cold, moderate, or hot).

```json
{
    "lat": 39.8097,
    "lon": -98.5556,
    "condition": "Clouds",
    "temperature": "cold"
}
```