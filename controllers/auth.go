package controllers

import (
	"horariosapp/database"
	"horariosapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Controller para renderizar la pagina de Inicio de session
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Controler con la logica de autenticacion
func Login(c *gin.Context) {

	var user models.User
	// Extraemos los inputs del formulario de login
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Verificamos si la base de datos existe
	if database.DB == nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Database not found"})
		return
	}

	// Verificamos si el usuario existe en la base de datos
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	// Verificamos si la contraseña es del usuario registrado
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	// Creamos la cookie de session del rol del usuario y otra con su nombre
	c.SetCookie("session", user.Role, 3600, "/", "", false, true)
	c.SetCookie("name", user.Name, 3600, "/", "", false, true)

	// Segun el rol del usuario lo redirigimos a una ruta determinada
	if user.Role == "admin" {
		c.Redirect(http.StatusFound, "/admin")
	} else {
		c.Redirect(http.StatusFound, "/login")
	}
}

// Controller con la logica de logout (eliminando las cookies)
func Logout(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
