package controllers

import (
	"horariosapp/database"
	"horariosapp/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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

	// Obtener los días de la semana para esta semana
	startDate, _ := time.Parse("2006-01-02", week.Start)
	days := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		days[i] = startDate.AddDate(0, 0, i)
	}

	// Obtener todos los trabajadores
	var workers []models.Worker
	if err := database.DB.Find(&workers).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{"error": err.Error()})
		return
	}

	// Formatear fechas en español
	formattedDays := make([]string, 7)
	for i, day := range days {
		formattedDays[i] = weekdays[day.Weekday()] + " " + strconv.Itoa(day.Day()) + " de " + months[day.Month()]
	}

	c.HTML(http.StatusOK, "week_detail.html", gin.H{
		"Week":      week,
		"Days":      formattedDays,
		"Intervals": Intervals,
		"Workers":   workers,
	})
}
