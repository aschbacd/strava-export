package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/antihax/optional"
	"github.com/aschbacd/strava-export/pkg/logger"
	swagger "github.com/aschbacd/strava-export/pkg/strava"
	"github.com/aschbacd/strava-export/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var (
	SHEETNAME = "Strava-Export"
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
		return activities[i].Date.Unix() < activities[j].Date.Unix()
	})

	// Create Excel file
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", SHEETNAME)

	// Set column widths
	for col, width := range map[string]float64{
		"A": 14, "B": 50, "C": 9, "D": 10, "E": 14, "F": 9, "G": 16, "H": 19, "I": 13, "J": 13, "K": 16, "L": 7, "M": 10, "N": 15} {
		if err := f.SetColWidth(SHEETNAME, col, col, width); err != nil {
			logger.Error(err.Error())
			utils.ReturnErrorPage(c)
			return
		}
	}

	// Define styles
	borderStyle := []excelize.Border{
		{
			Type:  "top",
			Color: "#000000",
			Style: 1,
		},
		{
			Type:  "right",
			Color: "#000000",
			Style: 1,
		},
		{
			Type:  "bottom",
			Color: "#000000",
			Style: 1,
		},
		{
			Type:  "left",
			Color: "#000000",
			Style: 1,
		},
	}

	// Format title
	if err := f.MergeCell(SHEETNAME, "A1", "N1"); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}
	if err := f.SetRowHeight(SHEETNAME, 1, 30); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}
	titleStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
			Size: 15,
		},
		Border: borderStyle,
	})
	if err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}
	if err := f.SetCellStyle(SHEETNAME, "A1", "N1", titleStyle); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Format header
	headerStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Border: borderStyle,
	})
	if err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}
	if err := f.SetCellStyle(SHEETNAME, "A2", "N2", headerStyle); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Format cells
	cellStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
		Border: borderStyle,
	})
	if err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}
	dateStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
		Border: borderStyle,
		NumFmt: 22,
	})
	if err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}
	durationStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
		Border: borderStyle,
		NumFmt: 21,
	})
	if err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}
	f.SetCellStyle(SHEETNAME, "A3", "A"+fmt.Sprint(len(activities)+2), dateStyle)
	f.SetCellStyle(SHEETNAME, "B3", "C"+fmt.Sprint(len(activities)+2), cellStyle)
	f.SetCellStyle(SHEETNAME, "D3", "D"+fmt.Sprint(len(activities)+2), durationStyle)
	f.SetCellStyle(SHEETNAME, "E3", "N"+fmt.Sprint(len(activities)+2), cellStyle)

	// Set title
	var title string
	firstDate := activities[0].DateLocal
	lastDate := activities[len(activities)-1].DateLocal

	if firstDate.Year() == lastDate.Year() && firstDate.Month() == lastDate.Month() {
		// If same year and same month -> Month - Year
		title = formatMonth(firstDate.Month()) + " - " + fmt.Sprint(firstDate.Year())
	} else if firstDate.Year() == lastDate.Year() {
		// If same year but different month -> Month1 - Month2 (Year)
		title = formatMonth(firstDate.Month()) + " - " + formatMonth(lastDate.Month()) + " (" + fmt.Sprint(firstDate.Year()) + ")"
	} else {
		// If different year and different month -> Month1 (Year1) - Month2 (Year2)
		title = formatMonth(firstDate.Month()) + " (" + fmt.Sprint(firstDate.Year()) + ") - " + formatMonth(lastDate.Month()) + " (" + fmt.Sprint(lastDate.Year()) + ")"
	}

	if err := f.SetCellValue(SHEETNAME, "A1", title); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Set Header (row 2)
	if err := setExcelValues(f, 2, []interface{}{
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
		"Fahrrad",
	}); err != nil {
		logger.Error(err.Error())
		utils.ReturnErrorPage(c)
		return
	}

	// Add activities to Excel file
	for i, activity := range activities {
		// Set values (row 3+)
		if err := setExcelValues(f, i+3, []interface{}{
			activity.DateLocal,
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
func setExcelValues(f *excelize.File, row int, items []interface{}) error {
	for i, item := range items {
		// Get cell names (e.g. A2, B7, C8, ...)
		axis, err := excelize.CoordinatesToCellName(i+1, row)
		if err != nil {
			return err
		}

		// Set value
		if err := f.SetCellValue(SHEETNAME, axis, item); err != nil {
			return err
		}
	}
	return nil
}

// formatMonth returns a month name for a given int
func formatMonth(month time.Month) string {
	months := []string{
		"Jänner", "Februar", "März", "April", "Mai", "Juni",
		"Juli", "August", "September", "Oktober", "November", "Dezember",
	}
	return months[month-1]
}
