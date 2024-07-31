package controllers

import (
	"horariosapp/database"
	"horariosapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Renderizamos la pagina de inicio de sesion
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Manejamos la logica de autenticacion
func Login(c *gin.Context) {

	var user models.User
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Verificamos que la base de datos existe
	if database.DB == nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Database connection error"})
		return
	}

	// Verificamos si el usuario existe en la base de datos
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	// Verificamos la contraseña del usuario
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	// Crear la cookie de sesión basada en el rol del usuario.
	c.SetCookie("session", user.Role, 3600, "/", "localhost", false, true)
	c.SetCookie("name", user.Name, 3600, "/", "localhost", false, true)

	// Redirigir basado en el rol.
	if user.Role == "admin" {
		c.Redirect(http.StatusFound, "/admin")
	} else {
		c.Redirect(http.StatusFound, "/login")
	}
}

// Logout borra las cookies de sesión y redirige al usuario a la página de inicio de sesión.
func Logout(c *gin.Context) {
	// Borrar las cookies de sesión.
	c.SetCookie("session", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/login")
}
