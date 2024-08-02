package controllers

import (
	"horariosapp/config"
	"horariosapp/database"
	"horariosapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Controlador para renderizar la página de inicio de sesión
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Controlador con la lógica de autenticación
func Login(c *gin.Context) {

	// Extraemos el usuario y la contraseña del formulario
	var user models.User
	username := c.PostForm("username")
	password := c.PostForm("password")

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	// Configuramos los claims del Token JWT
	claims := jwt.MapClaims{
		"name":     user.Name,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // El token expira en 1 hora
	}

	// Creamos el token y lo firmamos con la secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Error creating token"})
		return
	}
	c.SetCookie("token", tokenString, 3600, "/", "", false, true)

	// Redireccionamos según el rol
	switch user.Role {
	case "admin":
		c.Redirect(http.StatusFound, "/admin")
	case "worker":
		c.Redirect(http.StatusFound, "/worker/"+user.Username)
	case "store":
		c.Redirect(http.StatusFound, "/store/"+user.Username)
	default:
		c.Redirect(http.StatusFound, "/login")
	}

}

// Controlador con la lógica de logout (eliminando la cookie del token)
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", true, true) // Elimina la cookie
	c.Redirect(http.StatusFound, "/login")
}
