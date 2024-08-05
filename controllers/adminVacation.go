package controllers

import (
	"horariosapp/config"
	"horariosapp/database"
	"horariosapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VacationsPage(c *gin.Context) {

	// Verificamos el token
	claims, err := config.ValidateToken(c)
	if err != nil || claims["role"] != "admin" {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var workers []models.Worker
	if err := database.DB.Find(&workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var vacations []models.Vacation
	if err := database.DB.Preload("Worker").Find(&vacations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "vacations.html", gin.H{
		"Workers":   workers,
		"Vacations": vacations,
	})
}

// Controlador para crear vacaciones
func CreateVacation(c *gin.Context) {
	var input struct {
		WorkerID  uint   `form:"worker_id" binding:"required"`
		StartDate string `form:"start_date" binding:"required"`
		EndDate   string `form:"end_date" binding:"required"`
		Status    string `form:"status" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vacation := models.Vacation{
		WorkerID:  input.WorkerID,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		Status:    input.Status,
	}

	if err := database.DB.Create(&vacation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/vacations")
}

// Controlador para editar vacaciones
func UpdateVacation(c *gin.Context) {
	var input struct {
		ID        uint   `form:"id" binding:"required"`
		WorkerID  uint   `form:"worker_id" binding:"required"`
		StartDate string `form:"start_date" binding:"required"`
		EndDate   string `form:"end_date" binding:"required"`
		Status    string `form:"status" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var vacation models.Vacation
	if err := database.DB.Where("id = ?", input.ID).First(&vacation).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vacation not found"})
		return
	}

	database.DB.Model(&vacation).Updates(input)
	c.Redirect(http.StatusFound, "/admin/vacations")
}

// Controlador para eliminar vacaciones
func DeleteVacation(c *gin.Context) {
	var vacation models.Vacation
	if err := database.DB.Where("id = ?", c.Param("id")).First(&vacation).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vacation not found"})
		return
	}

	if err := database.DB.Delete(&vacation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vacation"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/vacations")
}
