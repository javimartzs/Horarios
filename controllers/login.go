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

// Controlador para abir la pagina de Login ("/" o "/login")
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// Controlador con la logica de autenticacion de los diferentes usuarios
func Login(c *gin.Context) {

	var user models.User
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Verificamos si el usuario esta en la base de datos
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}
	// Verificamos que esa contraseña es de ese usuario
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid username or password"})
		return
	}

	// Configuramos los claims del token JWT
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"name":     user.Name,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	// Creamos el token y lo firmamos con la key
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenStr, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Error creating auth token"})
		return
	}

	// Creamos la cookie asociada al token jwt
	c.SetCookie("token", tokenStr, 3600, "/", "", false, true)

	// Redireccionamos a la pagina segun el rol de la cookie
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
