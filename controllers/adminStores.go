package controllers

import (
	"horariosapp/config"
	"horariosapp/database"
	"horariosapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controlador para mostrar la pagina de tiendas
func StoresPage(c *gin.Context) {

	// Verificamos el token
	claims, err := config.ValidateToken(c)
	if err != nil || claims["role"] != "admin" {
		c.Redirect(http.StatusFound, "/login")
	}

	// Recuperamos el registro de tiendas
	var stores []models.Store
	if err := database.DB.Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Pasamos la lista a la plantilla HTML
	c.HTML(http.StatusOK, "adminStores.html", gin.H{
		"Stores": stores,
	})
}

// Estructura de datos para crear tiendas
type CreateStoreInput struct {
	Name   string `form:"name" binding:"required"`
	City   string `form:"city" binding:"required"`
	Phone  string `form:"phone" binding:"required"`
	Status string `form:"status" binding:"required"`
}

// Controller para crear una nueva tienda
func CreateStore(c *gin.Context) {

	var input CreateStoreInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store := models.Store{
		Name:   input.Name,
		City:   input.City,
		Phone:  input.Phone,
		Status: input.Status,
	}
	if err := database.DB.Create(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Redirigimos a stores de nuevo
	c.Redirect(http.StatusFound, "/admin/stores")

}

// Estructura para actualizar la informacion de las tiendas
type UpdateStoreInput struct {
	Name   string `form:"name"`
	City   string `form:"city"`
	Phone  string `form:"phone"`
	Status string `form:"status"`
}

// Controller para actualizar la informacion de las tiendas
func UpdateStore(c *gin.Context) {

	var store models.Store
	if err := database.DB.Where("id = ?", c.Param("id")).First(&store).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input UpdateStoreInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&store).Updates(input)
	c.Redirect(http.StatusFound, "/admin/stores")
}

// Controller para eliminar una tienda de la tabla
func DeleteStore(c *gin.Context) {

	// Buscamos la tienda
	var store models.Store
	if err := database.DB.Where("id = ?", c.Param("id")).First(&store).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Eliminamos la tienda
	if err := database.DB.Delete(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Redirigimos
	c.Redirect(http.StatusFound, "/admin/stores")

}
