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
