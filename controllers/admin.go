package controllers

import (
	"horariosapp/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Controller para mostrar el panel de administrador del Admin
func ShowAdminPage(c *gin.Context) {
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

	// Obtenemos el nombre del usuario desde los claims
	name := claims["name"].(string)

	// Obtenemos la fecha actual
	currentDate := time.Now().Format("Monday, 2 January 2006")

	// Renderizamos la plantilla con los datos
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"Name":        name,
		"CurrentDate": currentDate,
	})
}
