package controllers

import (
	"horariosapp/config"
	"horariosapp/database"
	"horariosapp/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Controller para mostrar la pagina de trabajadores
func WorkersPage(c *gin.Context) {

	// Verificamos el token
	claims, err := config.ValidateToken(c)
	if err != nil || claims["role"] != "admin" {
		c.Redirect(http.StatusFound, "/login")
	}

	// Recuperamos el registro de trabajadores
	var workers []models.Worker
	if err := database.DB.Find(&workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Recuperamos el registro de tiendas para el desplegable
	var stores []models.Store
	if err := database.DB.Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Crear sets para almacenar tiendas y estados únicos
	statusSet := make(map[string]struct{})
	for _, worker := range workers {
		statusSet[worker.Status] = struct{}{}
	}
	var statuses []string
	for status := range statusSet {
		statuses = append(statuses, status)
	}

	// Pasar las listas a la plantilla HTML
	c.HTML(http.StatusOK, "adminWorkers.html", gin.H{
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

// Controller para crear un nuevo trabajador en postgres
func CreateWorker(c *gin.Context) {

	// Validamos que los datos del formulario cumplan la estructura de CreateWorkerInpt
	var input CreateWorkerInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Crear un nuevo trabajador con los datos del formulario validado
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
	// Guardar el trabajador en la base de datos
	if err := database.DB.Create(&worker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Crear usuario de acceso para el nuevo trabajador
	username := strings.ToLower(input.Name) + strings.ToLower(input.Lastname)
	password := input.Identification
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Guardamos el nuevo trabajador como usuario en su base
	user := models.User{
		Username: username,
		Password: string(hashPassword),
		Role:     "worker",
		Name:     input.Name,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Al terminar la creacion me redirige a workers de nuevo
	c.Redirect(http.StatusFound, "/admin/workers")
}

// Estructura de datos para actualizar la informacion de los trabajadores
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

// Controlador para actualizar la informacion de los trabajadores
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

	// Buscamos al trabajador
	var worker models.Worker
	if err := database.DB.Where("id = ?", c.Param("id")).First(&worker).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Worker not Found"})
		return
	}

	// Eliminamos el trabajador
	if err := database.DB.Delete(&worker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Eliminamos el usuario asociado
	if err := database.DB.Where("username = ?", strings.ToLower(worker.Name+worker.Lastname)).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated user"})
		return
	}

	// Redirigimos despues de eliminar el usuario y el trabajador
	c.Redirect(http.StatusFound, "/admin/workers")
}
