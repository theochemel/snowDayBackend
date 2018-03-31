package main

type prediction struct {
  PredictionProbability float64 `json: PredictionProbability`
  PredictionBoolean bool `json: PredictionBoolean`
}

func makePrediction(forecast weatherConditions) prediction {
  var tomorrowPrediction prediction
  if forecast.Hourly.Data[0].Temperature < 32 && forecast.Hourly.Data[0].PrecipProbability > 0 {
    tomorrowPrediction.PredictionProbability = 100.0
    tomorrowPrediction.PredictionBoolean = true
  } else {
    tomorrowPrediction.PredictionProbability = 0.0
    tomorrowPrediction.PredictionBoolean = false
  }
  return tomorrowPrediction
}
