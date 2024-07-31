package controllers

import (
	"net/http"

	"horariosapp/database"
	"horariosapp/models"

	"github.com/gin-gonic/gin"
)

func ShowWorkersPage(c *gin.Context) {
	var workers []models.Worker
	if err := database.DB.Find(&workers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generar la lista de tiendas unicas a partir de la base de trabajadores
	storeSet := make(map[string]struct{})
	for _, worker := range workers {
		storeSet[worker.Store] = struct{}{}
	}
	var stores []string
	for store := range storeSet {
		stores = append(stores, store)
	}

	c.HTML(http.StatusOK, "workers.html", gin.H{
		"Workers": workers,
		"Stores":  stores,
	})
}

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

	// Redirigir a /admin/workers después de crear el trabajador
	c.Redirect(http.StatusFound, "/admin/workers")
}

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

func UpdateWorker(c *gin.Context) {
	var worker models.Worker
	if err := database.DB.Where("id = ?", c.Param("id")).First(&worker).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Worker not found!"})
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

func DeleteWorker(c *gin.Context) {
	var worker models.Worker
	if err := database.DB.Where("id = ?", c.Param("id")).First(&worker).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Worker not found!"})
		return
	}

	database.DB.Delete(&worker)

	c.Redirect(http.StatusFound, "/admin/workers")
}
