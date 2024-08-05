package controllers

import (
	"horariosapp/config"
	"horariosapp/database"
	"horariosapp/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func WeeksPage(c *gin.Context) {

	// Verificamos el token
	claims, err := config.ValidateToken(c)
	if err != nil || claims["role"] != "admin" {
		c.Redirect(http.StatusFound, "/login")
	}

	// Obtemeos el parametro del año del formulario, por defecto el actual
	yearParam := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, _ := strconv.Atoi(yearParam)

	// Recuperamos los registros de las semanas del año especificado
	var weeks []models.Week
	database.DB.Where("year = ?", year).Find(&weeks)

	years := []int{2024, 2025, 2026, 2027, 2028, 2029, 2030, 2031, 2032, 2033, 2034}
	currentDate := time.Now().Format("2006-01-02")

	c.HTML(http.StatusOK, "adminCalendar.html", gin.H{
		"Weeks":       weeks,
		"Year":        year,
		"Years":       years,
		"CurrentDate": currentDate,
	})

}
