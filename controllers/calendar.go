package controllers

import (
	"horariosapp/config"
	"horariosapp/database"
	"horariosapp/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Controller para cargar la pagina que muestra el listado de semanas
func ShowWeeksPage(c *gin.Context) {
	// Extraemos el token de la cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Verificamos el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Extraemos los claims del token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["role"] != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtenemos el parametro de año del formulario, por defecto el actual
	yearParam := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, _ := strconv.Atoi(yearParam) // Lo convertimos a un entero

	// Recuperamos los registros de semanas para el año especificado de la base postgres
	var weeks []models.Week
	database.DB.Where("year = ?", year).Find(&weeks)

	years := []int{2024, 2025, 2026, 2027, 2028, 2029, 2030, 2031, 2032, 2033, 2034}

	c.HTML(http.StatusOK, "calendar.html", gin.H{
		"Weeks":       weeks,
		"Year":        year,
		"Years":       years,
		"CurrentDate": time.Now().Format("2006-01-02"),
	})
}
