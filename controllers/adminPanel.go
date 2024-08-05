package controllers

import (
	"horariosapp/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ShowPageAdmin(c *gin.Context) {

	// Validacion de token
	claims, err := config.ValidateToken(c)
	if err != nil || claims["role"] != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Extraemos el nombre del usuario y la fecha local
	name := claims["name"].(string)
	date := time.Now().Format("Monday, 2 January 2006")

	// Renderizamos la plantilla con los datos
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"name": name,
		"date": date,
	})
}
