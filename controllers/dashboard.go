package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/optional"
	"github.com/aschbacd/strava-export/pkg/logger"
	swagger "github.com/aschbacd/strava-export/pkg/strava"
	"github.com/aschbacd/strava-export/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Activity represents a formatted activity
type Activity struct {
	Id               int64
	Date             string
	Name             string
	Distance         string
	Duration         string
	ElevationGain    string
	AverageSpeed     string
	MaxSpeed         string
	AverageCadence   string
	AverageHeartRate string
	MaxHeartRate     string
	AverageWatts     string
	MaxWatts         string
	Calories         string
	Kilojoules       string
	DeviceName       string
}

// ActivityDetails contains detailed information about the activity
type ActivityDetails struct {
	AverageCadence   float32 `json:"average_cadence"`
	AverageHeartRate float32 `json:"average_heartrate"`
	MaxHeartRate     float32 `json:"max_heartrate"`
	Calories         float32 `json:"calories"`
	DeviceName       string  `json:"device_name"`
}

// GetActivitiesPage returns the activities page
func GetActivitiesPage(c *gin.Context) {
	// Get page number (default = 1)
	pageNumber := 1

	if page := c.Query("page"); page != "" {
		// Convert string to int
		number, err := strconv.Atoi(page)
		if err != nil {
			logger.Error(err.Error())
			utils.ReturnErrorPage(c)
			return
		}
		pageNumber = number
	}

	// Create activities api config
	athleteActivityOpts := swagger.ActivitiesApiGetLoggedInAthleteActivitiesOpts{
		Page: optional.NewInt32(int32(pageNumber)),
	}

	// Set timestamps for activities api config
	if err := setAthleteActivitiesOpts(c, &athleteActivityOpts); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Get pagination links
	linkBefore := c.Request.URL.Query()
	linkBefore.Set("page", fmt.Sprint(pageNumber-1))
	linkAfter := c.Request.URL.Query()
	linkAfter.Set("page", fmt.Sprint(pageNumber+1))

	// Get activities (not detailed)
	activities, err := getActivities(c, athleteActivityOpts, false)
	if err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Return activities view
	c.HTML(http.StatusOK, "activities", gin.H{
		"activities": activities,
		"hasBefore":  pageNumber > 1,
		"hasAfter":   len(activities) == 30,
		"linkBefore": "?" + linkBefore.Encode(),
		"linkAfter":  "?" + linkAfter.Encode(),
		"from":       c.Query("from"),
		"to":         c.Query("to"),
	})
}

// setAthleteActivitiesOpts sets the timestamps for a given activities api config
func setAthleteActivitiesOpts(c *gin.Context, athleteActivityOpts *swagger.ActivitiesApiGetLoggedInAthleteActivitiesOpts) error {
	// From
	if after := c.Query("from"); after != "" {
		date, err := time.Parse("2006-01-02", after)
		if err != nil {
			return err
		}
		// Use last second of day before
		date = date.Add(-time.Second)
		athleteActivityOpts.After = optional.NewInt32(int32(date.Unix()))
	}

	// To
	if before := c.Query("to"); before != "" {
		date, err := time.Parse("2006-01-02", before)
		if err != nil {
			return err
		}
		// Use first second of day after
		date = date.Add(time.Hour * 24)
		athleteActivityOpts.Before = optional.NewInt32(int32(date.Unix()))
	}

	return nil
}

// getActivities generates a formatted list of activities
func getActivities(c *gin.Context, athleteActivityOpts swagger.ActivitiesApiGetLoggedInAthleteActivitiesOpts, detailed bool) ([]Activity, error) {
	// Get token source from authentication middleware
	tokenSource, exists := c.Get("tokenSource")
	if !exists {
		return nil, fmt.Errorf("client not passed by authentication middleware")
	}

	// Create new swagger client
	client := swagger.NewAPIClient(swagger.NewConfiguration())
	auth := context.WithValue(context.Background(), swagger.ContextOAuth2, tokenSource)

	// Create list
	activities := []Activity{}

	// Get activities from Strava
	stravaActivities, resp, _ := client.ActivitiesApi.GetLoggedInAthleteActivities(auth, &athleteActivityOpts)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get activity summary")
	}

	// Add activities to list
	for _, stravaActivity := range stravaActivities {
		// Parse duration from seconds
		duration, err := time.ParseDuration(fmt.Sprint(stravaActivity.MovingTime) + "s")
		if err != nil {
			return nil, err
		}

		// Create activity
		activity := Activity{
			Id:            stravaActivity.Id,
			Name:          stravaActivity.Name,
			Date:          stravaActivity.StartDateLocal.Format("02.01.2006 15:04:05"),
			Distance:      fmt.Sprintf("%.2f", stravaActivity.Distance/1000) + "km",
			Duration:      duration.String(),
			ElevationGain: fmt.Sprint(stravaActivity.TotalElevationGain) + "m",
			AverageSpeed:  fmt.Sprintf("%.2f", stravaActivity.AverageSpeed*3.6) + "km/h",
			MaxSpeed:      fmt.Sprintf("%.2f", stravaActivity.MaxSpeed*3.6) + "km/h",
			AverageWatts:  fmt.Sprint(stravaActivity.AverageWatts),
			MaxWatts:      fmt.Sprint(stravaActivity.MaxWatts),
			Kilojoules:    fmt.Sprintf("%.2f", stravaActivity.Kilojoules),
		}

		// Get activity details
		if detailed {
			// Get client from authentication middleware
			client, exists := c.Get("client")
			if !exists {
				return nil, fmt.Errorf("client not passed by authentication middleware")
			}

			// JSON response must be used instead of Object because some attributes are not supported
			resp, err := client.(*http.Client).Get("https://www.strava.com/api/v3/activities/" + fmt.Sprint(stravaActivity.Id))
			if err != nil {
				return nil, err
			}

			defer resp.Body.Close()

			stravaActivityDetails := ActivityDetails{}
			if err = json.NewDecoder(resp.Body).Decode(&stravaActivityDetails); err != nil {
				return nil, err
			}

			// Set details
			activity.AverageCadence = fmt.Sprintf("%.2f", stravaActivityDetails.AverageCadence)
			activity.AverageHeartRate = fmt.Sprintf("%.2f", stravaActivityDetails.AverageHeartRate)
			activity.MaxHeartRate = fmt.Sprintf("%.2f", stravaActivityDetails.MaxHeartRate)
			activity.Calories = fmt.Sprintf("%.2f", stravaActivityDetails.Calories)
			activity.DeviceName = stravaActivityDetails.DeviceName
		}

		// Add activity
		activities = append(activities, activity)
	}

	return activities, nil
}
