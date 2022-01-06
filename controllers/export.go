package controllers

import (
	"net/http"
	"sort"

	"github.com/antihax/optional"
	"github.com/aschbacd/strava-export/pkg/logger"
	swagger "github.com/aschbacd/strava-export/pkg/strava"
	"github.com/aschbacd/strava-export/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// ExportData exports an Excel report
func ExportData(c *gin.Context) {
	// Create activities api config
	athleteActivityOpts := swagger.ActivitiesApiGetLoggedInAthleteActivitiesOpts{
		PerPage: optional.NewInt32(60),
	}

	// Set timestamps for activities api config
	if err := setAthleteActivitiesOpts(c, &athleteActivityOpts); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Create Excel file
	f := excelize.NewFile()
	if err := setExcelValues(f, 1, []string{
		"Datum",
		"Name",
		"Strecke",
		"Zeit",
		"Höhenzunahme",
		"Kalorien",
		"Ø Geschwindigkeit",
		"Max. Geschwindigkeit",
		"Ø Trittfrequenz",
		"Ø Herzfrequenz",
		"Max. Herzfrequenz",
		"Ø Watt",
		"Max. Watt",
		"Rad",
	}); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Get activities (detailed)
	activities, rateLimitReached, errors := getActivities(c, athleteActivityOpts, true)
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

	// Sort activities
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Timestamp < activities[j].Timestamp
	})

	// Add activities to Excel file
	for i, activity := range activities {
		if err := setExcelValues(f, i+2, []string{
			activity.Date,
			activity.Name,
			activity.Distance,
			activity.Duration,
			activity.ElevationGain,
			activity.Calories,
			activity.AverageSpeed,
			activity.MaxSpeed,
			activity.AverageCadence,
			activity.AverageHeartRate,
			activity.MaxHeartRate,
			activity.AverageWatts,
			activity.MaxWatts,
			activity.GearName,
		}); err != nil {
			logger.Error(err.Error())
			utils.ReturnErrorPage(c)
			return
		}
	}

	// Set headers to make file downloadable
	c.Header("Content-Disposition", "attachment;filename=strava-export.xlsx")
	c.Header("File-Name", "strava-export.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")

	// Write file to gin's response writer
	f.Write(c.Writer)
}

// setExcelValues sets the values for a given row
func setExcelValues(f *excelize.File, row int, items []string) error {
	for i, item := range items {
		// Get cell names (e.g. A2, B7, C8, ...)
		axis, err := excelize.CoordinatesToCellName(i+1, row)
		if err != nil {
			return err
		}

		// Set value
		if err := f.SetCellValue("Sheet1", axis, item); err != nil {
			return err
		}
	}
	return nil
}
