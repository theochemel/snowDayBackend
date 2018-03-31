package main

import (
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "encoding/json"
  "fmt"
)

func startRouter() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", Index)
  router.HandleFunc("/{latitude},{longitude}", fetchForecastAndReturnPrediction)
  log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {

}

func fetchForecastAndReturnPrediction(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  lat := vars["latitude"]
  lon := vars["longitude"]

  forecast := pullDarkSkyForecast(lat, lon)
  fmt.Println(forecast.Currently.Summary)
  tomorrowPrediction := makePrediction(forecast)

  json.NewEncoder(w).Encode(tomorrowPrediction)

}
