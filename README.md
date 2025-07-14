# windyv1

The windyv1 is a simple http client application to retrieve the real-time weather from windyapi.com website for a given
airport code's latitude, longitude and apikey. The real-time response will be in the Json format.
It is up to the end user as how they want to process the response. The data could be saved for
additional processing or simply could be added in the database for other applications to process the data.

An example of invoking windy application is shown below.
resp, err := windy.GetWeather(53.1900, -112.2500, "valid api key")

Here's how to install it:

```
go install github.com/golangtrainingapp/windyv1@latest
```