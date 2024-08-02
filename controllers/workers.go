package controllers

import (
	"horariosapp/config"
	"horariosapp/database"
	"horariosapp/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func ShowWorkersPage(c *gin.Context) {
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

	// Recuperamos los registros de la tabla de postgres
	var workers []models.Worker
	if err := database.DB.Find(&workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generamos la lista de tiendas unicas para el filtro y la ordenacion de la tabla
	storeSet := make(map[string]struct{})
	statusSet := make(map[string]struct{})
	for _, worker := range workers {
		storeSet[worker.Store] = struct{}{}
		statusSet[worker.Status] = struct{}{}
	}
	var stores []string
	for store := range storeSet {
		stores = append(stores, store)
	}
	var statuses []string
	for status := range statusSet {
		statuses = append(statuses, status)
	}

	c.HTML(http.StatusOK, "workers.html", gin.H{
		"Workers":  workers,
		"Stores":   stores,
		"Statuses": statuses,
	})

}

// Estructura de datos para crear un nuevo trabajador
type CreateWorkerInput struct {
	Name           string `form:"name" binding:"required"`
	Lastname       string `form:"lastname" binding:"required"`
	Email          string `form:"email" binding:"required,email"`
	Identification string `form:"identification" binding:"required"`
	Cargo          string `form:"cargo" binding:"required"`
	Store          string `form:"store" binding:"required"`
	Status         string `form:"status" binding:"required"`
	PeriodoPrueba  string `form:"periodoPrueba" binding:"required"`
}

// Controlador para crear un nuevo trabajador en la base de postgres
func CreateWorker(c *gin.Context) {
	var input CreateWorkerInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	worker := models.Worker{
		Name:           input.Name,
		Lastname:       input.Lastname,
		Email:          input.Email,
		Identification: input.Identification,
		Cargo:          input.Cargo,
		Store:          input.Store,
		Status:         input.Status,
		PeriodoPrueba:  input.PeriodoPrueba,
	}

	if err := database.DB.Create(&worker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Crear usuario de acceso al trabajador
	username := strings.ToLower(input.Name + input.Lastname)
	password := input.Identification
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := models.User{
		Username: username,
		Password: string(hashpassword),
		Role:     "worker",
		Name:     input.Name,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Redirigir a /admin/workers despues de crear el trabajador
	c.Redirect(http.StatusFound, "/admin/workers")
}

// Estructura de datos para actualizar la informacion de un trabajador existente
type UpdateWorkerInput struct {
	Name           string `form:"name"`
	Lastname       string `form:"lastname"`
	Email          string `form:"email"`
	Identification string `form:"identification"`
	Cargo          string `form:"cargo"`
	Store          string `form:"store"`
	Status         string `form:"status"`
	PeriodoPrueba  string `form:"periodoPrueba"`
}

// Controlador para actualizar la info de un trabajador
func UpdateWorker(c *gin.Context) {
	var worker models.Worker
	if err := database.DB.Where("id = ?", c.Param("id")).First(&worker).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input UpdateWorkerInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&worker).Updates(input)
	c.Redirect(http.StatusFound, "/admin/workers")
}

// Controlador para eliminar a un trabajador de la tabla
func DeleteWorker(c *gin.Context) {
	// Buscar el trabajador por ID
	var worker models.Worker
	if err := database.DB.Where("id = ?", c.Param("id")).First(&worker).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Worker not found!"})
		return
	}

	// Eliminar el trabajador
	if err := database.DB.Delete(&worker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete worker"})
		return
	}

	// Eliminar el usuario asociado
	if err := database.DB.Where("username = ?", strings.ToLower(worker.Name+worker.Lastname)).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated user"})
		return
	}

	// Redirigir a /admin/workers después de eliminar el trabajador y el usuario
	c.Redirect(http.StatusFound, "/admin/workers")
}
