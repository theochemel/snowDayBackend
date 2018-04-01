package main

import (
	"fmt"
)


type stormInfo struct {
  StartTime int `json: "startTime"`
  EndTime int `json: "endTime"`
  PeakIntensity float64 `json: "peakIntensity"`
  PeakIntesityTime int `json: "peakIntesityTime"`
  TotalAccumulation float64 `json: "totalAccumulation"`
  AverageProbability float64 `json: "averageProbability"`
}

func formatForecast(forecast weatherConditions) stormInfo {
  storm := stormInfo{}
  stormStarted := false
  totalProbability := 0.0
  numSamples := 0.0
  runningPeakIntensity := 0.0
  runningPeakIntensityTime := 0
  fmt.Println(len(forecast.Hourly.Data))

  Loop:
    for _, point := range forecast.Hourly.Data {
      fmt.Println(point)
      if point.PrecipProbability > 0.1 && point.PrecipType == "snow" {
        if stormStarted == false {
          storm.StartTime = point.Time
          stormStarted = true
          fmt.Println(storm.StartTime)
        }
        storm.TotalAccumulation += point.PrecipAccumulation
        totalProbability += point.PrecipProbability
        numSamples++

        if point.PrecipIntensity > runningPeakIntensity {
          runningPeakIntensity = point.PrecipIntensity
          runningPeakIntensityTime = point.Time
        }

      } else {
        if stormStarted == true {
          storm.EndTime = point.Time
          fmt.Println(point.Time)
          break Loop
        }
      }
    }

  fmt.Println(storm.TotalAccumulation)
  storm.PeakIntensity = runningPeakIntensity
  storm.PeakIntesityTime = runningPeakIntensityTime
  storm.AverageProbability = totalProbability / numSamples
  fmt.Println(storm)
  return storm
}
