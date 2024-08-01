package controllers

import (
	"encoding/json"
	"fmt"
	"horariosapp/database"
	"horariosapp/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Visualizacion de la lista de semanas por años
func ShowWeeksPage(c *gin.Context) {
	yearParam := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, _ := strconv.Atoi(yearParam)

	var weeks []models.Week
	database.DB.Where("year = ?", year).Find(&weeks)

	years := []int{2024, 2025, 2026, 2027, 2028, 2029, 2030, 2031, 2032, 2033, 2034}

	c.HTML(http.StatusOK, "horarios.html", gin.H{
		"Weeks": weeks,
		"Year":  year,
		"Years": years,
	})

}

var Intervals = []string{
	"08:00", "08:30", "09:00", "09:30", "10:00", "10:30", "11:00", "11:30",
	"12:00", "12:30", "13:00", "13:30", "14:00", "14:30", "15:00", "15:30",
	"16:00", "16:30", "17:00", "17:30", "18:00", "18:30", "19:00", "19:30",
	"20:00", "20:30", "21:00", "21:30", "22:00",
}

// Mapas para traducción
var weekdays = map[time.Weekday]string{
	time.Sunday:    "Domingo",
	time.Monday:    "Lunes",
	time.Tuesday:   "Martes",
	time.Wednesday: "Miércoles",
	time.Thursday:  "Jueves",
	time.Friday:    "Viernes",
	time.Saturday:  "Sábado",
}

var months = map[time.Month]string{
	time.January:   "Enero",
	time.February:  "Febrero",
	time.March:     "Marzo",
	time.April:     "Abril",
	time.May:       "Mayo",
	time.June:      "Junio",
	time.July:      "Julio",
	time.August:    "Agosto",
	time.September: "Septiembre",
	time.October:   "Octubre",
	time.November:  "Noviembre",
	time.December:  "Diciembre",
}

func ShowWeekPage(c *gin.Context) {
	weekID := c.Param("weekID")

	var week models.Week
	if err := database.DB.First(&week, "week_id = ?", weekID).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{"error": "Week not found"})
		return
	}

	startDate, _ := time.Parse("2006-01-02", week.Start)
	days := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		days[i] = startDate.AddDate(0, 0, i)
	}

	var workers []models.Worker
	if err := database.DB.Find(&workers).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	var scheduleEntries []models.ScheduleEntry
	database.DB.Where("week_id = ?", weekID).Find(&scheduleEntries)

	entriesByDayAndWorker := make(map[int]map[uint]map[string]string)
	cellColors := make(map[string]string)
	for _, entry := range scheduleEntries {
		if entriesByDayAndWorker[entry.DayIndex] == nil {
			entriesByDayAndWorker[entry.DayIndex] = make(map[uint]map[string]string)
		}
		if entriesByDayAndWorker[entry.DayIndex][entry.WorkerID] == nil {
			entriesByDayAndWorker[entry.DayIndex][entry.WorkerID] = make(map[string]string)
		}
		entriesByDayAndWorker[entry.DayIndex][entry.WorkerID][entry.Interval] = entry.Color
		cellColors[fmt.Sprintf("%d-%s-%d", entry.WorkerID, entry.Interval, entry.DayIndex)] = entry.Color
	}

	formattedDays := make([]string, 7)
	for i, day := range days {
		formattedDays[i] = weekdays[day.Weekday()] + " " + strconv.Itoa(day.Day()) + " de " + months[day.Month()]
	}

	// Convertir cellColors a JSON escapado
	cellColorsJSON, err := json.Marshal(cellColors)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": "Error al convertir colores a JSON"})
		return
	}

	c.HTML(http.StatusOK, "week_detail.html", gin.H{
		"Week":               week,
		"Days":               formattedDays,
		"Intervals":          Intervals,
		"Workers":            workers,
		"EntriesByDayWorker": entriesByDayAndWorker,
		"CellColors":         template.JS(cellColorsJSON), // Inyectar como cadena JSON escapada
	})
}

// ----------------------------------------------------------------------------------------------------------
func SaveSchedule(c *gin.Context) {
	weekID := c.Param("weekID")

	var colors map[string]string
	if err := c.BindJSON(&colors); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid data"})
		return
	}

	for key, color := range colors {
		parts := strings.Split(key, "-")
		if len(parts) != 3 {
			continue
		}

		workerID, _ := strconv.Atoi(parts[0])
		interval := parts[1]
		dayIndex, _ := strconv.Atoi(parts[2])

		var entry models.ScheduleEntry
		err := database.DB.Where("week_id = ? AND worker_id = ? AND interval = ? AND day_index = ?", weekID, workerID, interval, dayIndex).First(&entry).Error

		if err != nil {
			// Si no existe la entrada, crear una nueva
			entry = models.ScheduleEntry{
				WeekID:   weekID,
				WorkerID: uint(workerID),
				Interval: interval,
				DayIndex: dayIndex,
				Color:    color,
			}
			database.DB.Create(&entry)
		} else {
			entry.Color = color
			database.DB.Save(&entry)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})

}
