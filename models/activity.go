package models

import (
	"time"
)

type Activity struct {
	Id               int64
	Date             time.Time
	DateLocal        time.Time
	Name             string
	Distance         float64 // [km]
	Duration         time.Duration
	ElevationGain    float64 // [m]
	AverageSpeed     float64 // [km/h]
	MaxSpeed         float64 // [km/h]
	AverageWatts     float64
	MaxWatts         int32
	Kilojoules       float64
	AverageCadence   float64
	AverageHeartRate float64
	MaxHeartRate     float64
	Calories         float64
	GearName         string
}

type ActivityDetails struct {
	AverageCadence   float64 `json:"average_cadence"`
	AverageHeartRate float64 `json:"average_heartrate"`
	MaxHeartRate     float64 `json:"max_heartrate"`
	Calories         float64 `json:"calories"`
	Gear             struct {
		Name string `json:"name"`
	} `json:"gear"`
}

func (a *Activity) GetDateString() string {
	return a.DateLocal.Format("02.01.2006")
}

func (a *Activity) GetTimeString() string {
	return a.DateLocal.Format("15:04:05")
}
