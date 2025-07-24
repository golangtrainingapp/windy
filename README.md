# windy

The windy application is a simple http client application to retrieve the real-time weather from windyapi.com website for a given
airport code's latitude, longitude and apikey. The real-time response will be in the Json format.
It is up to the end user as how they want to process the response. The data could be saved for
additional processing or simply could be added in the database for other applications to process the data.

Here's how to install it:

```
go install github.com/golangtrainingapp/windy/cmd/windy@latest
```

Type windy from the terminal. The application should retrieve the real-time json weather response for a sample latitude and longitude from windyapi.com