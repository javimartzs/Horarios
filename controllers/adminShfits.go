package controllers

import (
	"encoding/json"
	"fmt"
	"horariosapp/config"
	"horariosapp/database"
	"horariosapp/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Intervalos horarios de las tablas diarias
var Intervals = []string{
	"07:00", "07:30", "08:00", "08:30", "09:00", "09:30", "10:00", "10:30",
	"11:00", "11:30", "12:00", "12:30", "13:00", "13:30", "14:00", "14:30",
	"15:00", "15:30", "16:00", "16:30", "17:00", "17:30", "18:00", "18:30",
	"19:00", "19:30", "20:00", "20:30", "21:00", "21:30", "22:00",
}

// Mapas para traducción de días y meses
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

// ShowWeekPage maneja la visualización de la página de una semana específica
func ShowWeekPage(c *gin.Context) {

	// Verificamos el token
	claims, err := config.ValidateToken(c)
	if err != nil || claims["role"] != "admin" {
		c.Redirect(http.StatusFound, "/login")
	}

	// Obtiene el parámetro 'weekID' de la URL
	weekID := c.Param("weekID")

	// Busca la semana en la base de datos
	var week models.Week
	if err := database.DB.First(&week, "week_id = ?", weekID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Week not found"})
		return
	}

	// Calcula las fechas de los 7 días de la semana
	startDate, _ := time.Parse("2006-01-02", week.Start)
	days := make([]time.Time, 7)
	for i := 0; i < 7; i++ {
		days[i] = startDate.AddDate(0, 0, i)
	}

	// Recupera los trabajadores de la base de datos
	var workers []models.Worker
	if err := database.DB.Where("status = ?", "Alta").Find(&workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Recupera las entradas del horario para la semana específica
	var scheduleEntries []models.ScheduleEntry
	database.DB.Where("week_id = ?", weekID).Find(&scheduleEntries)

	// Organiza las entradas de horario por día y trabajador
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

	// Formatea las fechas para mostrar en la vista
	formattedDays := make([]string, 7)
	for i, day := range days {
		formattedDays[i] = weekdays[day.Weekday()] + " " + strconv.Itoa(day.Day()) + " de " + months[day.Month()]
	}

	// Recupera los totales de horas trabajadas
	var workerTotals []models.WorkerHours
	database.DB.Where("week_id = ?", weekID).Find(&workerTotals)

	// Organiza los totales de horas por trabajador y día
	totalsByWorkerAndDay := make(map[uint]map[int]float64)
	for _, total := range workerTotals {
		if totalsByWorkerAndDay[total.WorkerID] == nil {
			totalsByWorkerAndDay[total.WorkerID] = make(map[int]float64)
		}
		totalsByWorkerAndDay[total.WorkerID][total.DayIndex] = total.TotalHours
	}

	// Crea un resumen semanal para cada trabajador
	weeklySummaries := []map[string]interface{}{}
	for _, worker := range workers {
		totalHours := 0.0
		for _, dayTotals := range totalsByWorkerAndDay[worker.ID] {
			totalHours += dayTotals
		}
		weeklySummaries = append(weeklySummaries, map[string]interface{}{
			"WorkerName": worker.Name,
			"Store":      worker.Store,
			"TotalHours": totalHours,
		})
	}

	// Obtiene la lista de tiendas
	storeSet := make(map[string]struct{})
	for _, worker := range workers {
		storeSet[worker.Store] = struct{}{}
	}

	stores := make([]string, 0, len(storeSet))
	for store := range storeSet {
		stores = append(stores, store)
	}

	// Convierte cellColors y totalsByWorkerAndDay a JSON
	cellColorsJSON, err := json.Marshal(cellColors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al convertir colores a JSON"})
		return
	}

	totalsJSON, err := json.Marshal(totalsByWorkerAndDay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al convertir totales a JSON"})
		return
	}

	// Renderiza la vista con los datos obtenidos
	c.HTML(http.StatusOK, "adminShift.html", gin.H{
		"Week":               week,
		"Days":               formattedDays,
		"Intervals":          Intervals,
		"Workers":            workers,
		"EntriesByDayWorker": entriesByDayAndWorker,
		"CellColors":         template.JS(cellColorsJSON), // Inyecta JSON escapado en la plantilla
		"WorkerTotals":       template.JS(totalsJSON),     // Inyecta JSON escapado en la plantilla
		"Stores":             stores,
		"WeeklySummaries":    weeklySummaries, // Resumen semanal de horas
	})
}

// ----------------------------------------------------------------------------------------------------------

// SaveSchedule maneja el guardado de las entradas de horario y los totales en la base de datos
func SaveSchedule(c *gin.Context) {
	weekIDStr := c.Param("weekID")
	weekID, err := strconv.ParseUint(weekIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid weekID"})
		return
	}

	var requestData struct {
		Colors map[string]string             `json:"colors"`
		Totals map[string]map[string]float64 `json:"totals"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid data"})
		return
	}

	colors := requestData.Colors
	totals := requestData.Totals

	// Almacenar entradas de colores
	for key, color := range colors {
		parts := strings.Split(key, "-")
		if len(parts) != 3 {
			continue
		}

		workerID, _ := strconv.Atoi(parts[0])
		interval := parts[1]
		dayIndex, _ := strconv.Atoi(parts[2])

		var entry models.ScheduleEntry
		err := database.DB.Where("week_id = ? AND worker_id = ? AND interval = ? AND day_index = ?", uint(weekID), workerID, interval, dayIndex).First(&entry).Error

		if err != nil {
			entry = models.ScheduleEntry{
				WeekID:   uint(weekID),
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

	// Guardar totales en la base de datos
	for workerID, days := range totals {
		workerIDInt, _ := strconv.Atoi(workerID)
		for dayIndexStr, totalHours := range days {
			dayIndex, _ := strconv.Atoi(dayIndexStr)
			var workerTotal models.WorkerHours
			err := database.DB.Where("worker_id = ? AND week_id = ? AND day_index = ?", uint(workerIDInt), uint(weekID), dayIndex).First(&workerTotal).Error

			if err != nil {
				workerTotal = models.WorkerHours{
					WorkerID:   uint(workerIDInt),
					WeekID:     uint(weekID),
					DayIndex:   dayIndex,
					TotalHours: totalHours,
				}
				database.DB.Create(&workerTotal)
			} else {
				workerTotal.TotalHours = totalHours
				database.DB.Save(&workerTotal)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
