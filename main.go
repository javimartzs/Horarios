package main

import (
	"horariosapp/config"
	"horariosapp/controllers"
	"horariosapp/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Importamos las variables de entorno
	config.Init()
	// Iniciamos la base de datos
	database.ConnectDB()
	// Creamos el router de Gin
	r := gin.Default()
	// Importamos las plantillas HTML
	r.LoadHTMLGlob("templates/*")

	// Rutas y controllers
	r.GET("", controllers.ShowLoginPage)
	r.GET("/login", controllers.ShowLoginPage)

	r.Run(":8080")

}
