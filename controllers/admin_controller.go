package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Controller para mostrar el panel de administrador del Admin
func ShowAdminPage(c *gin.Context) {

	// Comprobamos el rol que tiene
	role, err := c.Cookie("session")
	if err != nil || role != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtenemos el nombre del usuario de la cookie
	name, err := c.Cookie("name")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtenemos la fecha actual
	currentDate := time.Now().Format("Monday, 2 January 2006")

	// Renderizamos la plantilla con los datos
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"Name":        name,
		"CurrentDate": currentDate,
	})
}
