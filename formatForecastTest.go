package main

type stormInfo struct {
	TimeForecastRecieved int `json: "timeForecastRecieved"`
	StormPossible bool `json: "stormPossible"`
  StartTime int `json: "startTime,omitempty"`
  EndTime int `json: "endTime,omitempty"`
  PeakIntensity float64 `json: "peakIntensity,omitempty"`
  PeakIntesityTime int `json: "peakIntesityTime,omitempty"`
  TotalAccumulation float64 `json: "totalAccumulation,omitempty"`
  AverageProbability float64 `json: "averageProbability,omitempty"`
}

func formatForecast(forecast weatherConditions) stormInfo {
  storm := stormInfo{}
	storm.StormPossible = false
	storm.TimeForecastRecieved = forecast.Hourly.Data[0].Time
  stormStarted := false
  totalProbability := 0.0
  numSamples := 0.0
  runningPeakIntensity := 0.0
  runningPeakIntensityTime := 0

  Loop:
    for _, point := range forecast.Hourly.Data {
      if point.PrecipProbability > 0.1 && point.PrecipType == "snow" {
        if stormStarted == false {
					storm.StormPossible = true
          storm.StartTime = point.Time
          stormStarted = true
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
          break Loop
        }
      }
		}
  storm.PeakIntensity = runningPeakIntensity
  storm.PeakIntesityTime = runningPeakIntensityTime
  storm.AverageProbability = totalProbability / numSamples
	if stormStarted == false {
		storm.AverageProbability = 0.0
	}
  return storm
}
