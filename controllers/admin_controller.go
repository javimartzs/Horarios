package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ShowAdminDashboard muestra el panel de administración solo si el usuario tiene el rol de administrador.
func ShowAdminPage(c *gin.Context) {
	role, err := c.Cookie("session")
	if err != nil || role != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener el nombre de usuario de la cookie o sesión
	name, err := c.Cookie("name")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener la fecha actual
	currentDate := time.Now().Format("Monday, 02 January 2006")

	c.HTML(http.StatusOK, "admin.html", gin.H{
		"name":        name,
		"currentDate": currentDate,
	})
}
