package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)



type weatherConditions struct {

  Currently struct {
    Time int `json: "time"`
    Temperature float64 `json: "temperature"`
    Summary string `json: "summary"`
    PrecipProbability float64 `json: "precipProbabilty"`
    PrecipIntensity float64 `json: "precipIntensity"`
  } `json: "currently"`

  Hourly struct {
    Summary string `json: "summary"`
    Data []struct {
      Time int `json: "time"`
      Temperature float64 `json: "temperature"`
      Summary string `json: "summary"`
      PrecipProbability float64 `json: "precipProbabilty"`
      PrecipIntensity float64 `json: "precipIntensity"`
      PrecipAccumulation float64 `json: precipAccumulation`
      PrecipType string `json: precipType`
    } `json: "data"`
  } `json: "hourly"`
}


func pullDarkSkyForecast(lat string, lon string) weatherConditions {
  url := "https://api.darksky.net/forecast/" + darkSkyAPIKey + "/" + lat +", " + lon

  darkSkyClient := http.Client {
    Timeout: time.Second * 2,
  }

  req, err := http.NewRequest(http.MethodGet, url, nil)
  if err != nil {
    log.Fatal(err)
  }

  res, getErr := darkSkyClient.Do(req)
  if getErr != nil {
		log.Fatal(getErr)
	}
  body, readErr := ioutil.ReadAll(res.Body)
  if readErr != nil {
    log.Fatal(readErr)
  }

	results := weatherConditions{}
	jsonErr := json.Unmarshal(body, &results)
	if jsonErr != nil {
    log.Fatal(jsonErr)
    fmt.Println(jsonErr)
  }
	return results
}
