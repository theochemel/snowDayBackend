package main

import (
 "log"
 "encoding/json"
 "github.com/aws/aws-lambda-go/events"
 "github.com/aws/aws-lambda-go/lambda"
 "net/http"
 "time"
 "io/ioutil"
 "os"
)


type weatherConditions struct {

  Currently struct {
    Time int `json:"time"`
    Temperature float64 `json:"temperature"`
    Summary string `json:"summary"`
    PrecipProbability float64 `json:"precipProbabilty"`
    PrecipIntensity float64 `json:"precipIntensity"`
  } `json:"currently"`

  Hourly struct {
    Summary string `json:"summary"`
    Data []struct {
      Time int `json:"time"`
      Temperature float64 `json:"temperature"`
      Summary string `json:"summary"`
      PrecipProbability float64 `json:"precipProbabilty"`
      PrecipIntensity float64 `json:"precipIntensity"`
      PrecipAccumulation float64 `json:precipAccumulation`
      PrecipType string `json:precipType`
    } `json:"data"`
  } `json:"hourly"`
}


type stormInfo struct {
	TimeForecastRecieved int `json:"timeForecastRecieved"`
	StormPossible bool `json:"stormPossible"`
  StartTime int `json:"startTime"`
  EndTime int `json:"endTime"`
  PeakIntensity float64 `json:"peakIntensity"`
  PeakIntensityTime int `json:"peakIntensityTime"`
  TotalAccumulation float64 `json:"totalAccumulation"`
  AverageProbability float64 `json:"averageProbability"`
  PrecipInfo[]float64 `json:"precipInfo"`
}

var darkSkyAPIKey = os.Getenv("DARKSKYAPIKEY")
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

  }
	return results
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
      storm.PrecipInfo = append(storm.PrecipInfo, point.PrecipIntensity)
		}
  storm.PeakIntensity = runningPeakIntensity
  storm.PeakIntensityTime = runningPeakIntensityTime
  storm.AverageProbability = totalProbability / numSamples
	if stormStarted == false {
		storm.AverageProbability = 0.0
	}
  return storm
}


// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

 // stdout and stderr are sent to AWS CloudWatch Logs
 log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

 // If no name is provided in the HTTP request body, throw an error
 latitude := request.QueryStringParameters["longitude"]
 longitude := request.QueryStringParameters["longitude"]
 forecast := pullDarkSkyForecast(latitude, longitude)
 formattedForecast := formatForecast(forecast)
 b, _:= json.Marshal(formattedForecast)
 return events.APIGatewayProxyResponse{
  Body: string(b),
  StatusCode: 200,
 }, nil

}

func main() {
 lambda.Start(Handler)
}
