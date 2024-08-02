package main

import (
	"horariosapp/config"
	"horariosapp/controllers"
	"horariosapp/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	//Cargamos las variables de entorno
	config.Init()
	// Iniciamos la base de datos
	database.ConnectDB()

	// Configuramos GIN
	r := gin.Default()

	// Sirve archivos estáticos desde la carpeta 'css'
	r.Static("/css", "./css")

	// Importamos los templates HTML
	r.LoadHTMLGlob("templates/*")

	// Rutas para la autenticacion
	r.GET("", controllers.ShowLoginPage)
	r.GET("/login", controllers.ShowLoginPage)
	r.POST("/login", controllers.Login)
	r.POST("/logout", controllers.Logout)
	r.GET("/logout", controllers.Logout)

	// Grupo de rutas admin
	admin := r.Group("/admin")
	{
		admin.GET("", controllers.ShowAdminPage)

		// Rutas de trabajadores bajo el grupo admin
		admin.GET("/workers", controllers.ShowWorkersPage)
		admin.POST("/workers/create", controllers.CreateWorker)
		admin.POST("/workers/update/:id", controllers.UpdateWorker)
		admin.POST("/workers/delete/:id", controllers.DeleteWorker)

		// Rutas para el calendario y las semanas bajo el grupo admin
		admin.GET("/calendar", controllers.ShowWeeksPage)
		admin.GET("/calendar/:weekID", controllers.ShowWeekPage)
		admin.POST("/calendar/:weekID/save", controllers.SaveSchedule)
	}
	// Iniciamos el servidor
	r.Run(":8080")
	log.Print("Server running on port :8080")
}
