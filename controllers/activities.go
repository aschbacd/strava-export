package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/antihax/optional"
	"github.com/aschbacd/strava-export/models"
	"github.com/aschbacd/strava-export/pkg/logger"
	swagger "github.com/aschbacd/strava-export/pkg/strava"
	"github.com/aschbacd/strava-export/pkg/utils"
	"github.com/gin-gonic/gin"
)

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
	activities, rateLimitReached, errors := getActivities(c, athleteActivityOpts, false)
	if len(errors) > 0 || rateLimitReached {
		// Log all errors
		for _, err := range errors {
			logger.Error(err.Error())
		}

		// Check if rate limit reached
		if rateLimitReached {
			c.Redirect(http.StatusFound, "/rate-limit")
		} else {
			utils.ReturnErrorPage(c)
		}

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
func getActivities(c *gin.Context, athleteActivityOpts swagger.ActivitiesApiGetLoggedInAthleteActivitiesOpts, detailed bool) ([]models.Activity, bool, []error) {
	// Get token source from authentication middleware
	tokenSource, exists := c.Get("tokenSource")
	if !exists {
		return nil, false, []error{fmt.Errorf("client not passed by authentication middleware")}
	}

	// Create new swagger client
	client := swagger.NewAPIClient(swagger.NewConfiguration())
	auth := context.WithValue(context.Background(), swagger.ContextOAuth2, tokenSource)

	// Create list
	activities := []models.Activity{}

	// Get activities from Strava
	stravaActivities, resp, _ := client.ActivitiesApi.GetLoggedInAthleteActivities(auth, &athleteActivityOpts)
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, true, []error{fmt.Errorf("rate limit reached")}
	} else if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return nil, false, []error{fmt.Errorf("failed to get activity summary")}
	}

	// Create channels for details
	channelActivities := make(chan models.Activity, 100)
	channelErrors := make(chan error, 100)

	var wg sync.WaitGroup
	errors := []error{}

	// Add activities to list
	for _, stravaActivity := range stravaActivities {
		// Convert int into duration
		duration, err := time.ParseDuration(fmt.Sprint(stravaActivity.MovingTime) + "s")
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to convert int into duration"))
		}

		// Create activity
		activity := models.Activity{
			Id:            stravaActivity.Id,
			Name:          stravaActivity.Name,
			Date:          stravaActivity.StartDate,
			DateLocal:     stravaActivity.StartDateLocal,
			Distance:      math.Round(float64(stravaActivity.Distance/10)) / 100,
			Duration:      duration,
			ElevationGain: math.Round(float64(stravaActivity.TotalElevationGain*100)) / 100,
			AverageSpeed:  math.Round(float64(stravaActivity.AverageSpeed*360)) / 100,
			MaxSpeed:      math.Round(float64(stravaActivity.MaxSpeed*360)) / 100,
			AverageWatts:  math.Round(float64(stravaActivity.AverageWatts*100)) / 100,
			MaxWatts:      stravaActivity.MaxWatts,
			Kilojoules:    math.Round(float64(stravaActivity.Kilojoules*100)) / 100,
		}

		// Get activity details
		if detailed {
			wg.Add(1)
			go getActivityDetails(c, activity, channelActivities, channelErrors, &wg)
		} else {
			// Add activity
			activities = append(activities, activity)
		}
	}

	// Wait for goroutines to finish
	wg.Wait()
	close(channelActivities)
	close(channelErrors)

	// Add activities
	for detailedActivity := range channelActivities {
		activities = append(activities, detailedActivity)
	}

	// Add errors
	for err := range channelErrors {
		errors = append(errors, err)
	}

	return activities, false, errors
}

// getActivityDetails fetches the details for an activity and pushes it to a channel
func getActivityDetails(c *gin.Context, activity models.Activity, activities chan<- models.Activity, errors chan<- error, wg *sync.WaitGroup) {
	// Get client from authentication middleware
	client, exists := c.Get("client")
	if !exists {
		errors <- fmt.Errorf("client not passed by authentication middleware")
		wg.Done()
		return
	}

	// JSON response must be used instead of Object because some attributes are not supported
	resp, err := client.(*http.Client).Get("https://www.strava.com/api/v3/activities/" + fmt.Sprint(activity.Id))
	if err != nil {
		errors <- err
		wg.Done()
		return
	}

	defer resp.Body.Close()
	var stravaActivityDetails models.ActivityDetails
	if err = json.NewDecoder(resp.Body).Decode(&stravaActivityDetails); err != nil {
		errors <- err
		wg.Done()
		return
	}

	// Set activity details
	activity.AverageCadence = math.Round(float64(stravaActivityDetails.AverageCadence*100)) / 100
	activity.AverageHeartRate = math.Round(float64(stravaActivityDetails.AverageHeartRate*100)) / 100
	activity.MaxHeartRate = math.Round(float64(stravaActivityDetails.MaxHeartRate*100)) / 100
	activity.Calories = math.Round(float64(stravaActivityDetails.Calories*100)) / 100
	activity.GearName = stravaActivityDetails.Gear.Name

	// Push activity to activities
	activities <- activity
	wg.Done()
}
